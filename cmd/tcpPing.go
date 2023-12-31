/*
Copyright © 2024 Mikael Schultz <mikael@conf-t.se>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bitcanon/iptool/ip"
	"github.com/bitcanon/iptool/tcp"
	"github.com/bitcanon/iptool/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var csvFlagError = errors.New("the --csv flag requires the --output-file flag to be set")

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping <destination> [port]",
	Short: "Send a sequence of TCP pings to a host",
	Long: `Send a sequence of TCP pings to a host.

The TCP ping command sends SYN packets to a host and
prints the response time, until the user presses Ctrl-C.

If no port is specified, the default port 443 is used.

Example:
  iptool tcp ping 1.0.0.1
  iptool tcp ping 1.0.0.1 443
  iptool tcp ping 1.0.0.1:53 --timeout 500`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check that the user provided one or two arguments
		if len(args) < 1 || len(args) > 2 {
			return errors.New("invalid number of arguments")
		}

		// Check if the user used the format host:port
		if strings.Contains(args[0], ":") {
			// Split the host and port
			hostPort := strings.Split(args[0], ":")
			args[0] = hostPort[0]
			args = append(args, hostPort[1])
		}

		// Parse the host
		host := args[0]

		// Parse the port
		port := 443
		if len(args) == 2 {
			// Convert the port to an integer
			p, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			// Check that the port is valid
			if p < 1 || p > 65535 {
				return errors.New("invalid port number, must be between 1 and 65535")
			}

			// Set the port
			port = p
		}

		return tcpPingAction(os.Stdout, host, port)
	},
}

func tcpPingAction(out io.Writer, host string, port int) error {
	// Define the delay duration
	delay := viper.GetDuration("tcp.ping.delay") * time.Millisecond

	// Define the number of packets to send
	count := viper.GetInt("tcp.ping.count")

	// If the --csv flag is set and --output-file is not set, return an error
	if viper.GetBool("tcp.ping.csv") && !viper.IsSet("tcp.ping.output-file") {
		return csvFlagError
	}

	// Resolve the IP address of the destination
	ip, err := ip.ResolveIP(host)
	if err != nil {
		return err
	}

	// Create a channel to receive interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Packet counters
	packetsSent := 0
	packetsReceived := 0

	// Response times
	minResponseTime := time.Duration(0)
	maxResponseTime := time.Duration(0)
	avgResponseTime := time.Duration(0)
	mdevResponseTime := time.Duration(0)
	totResponseTime := time.Duration(0)
	totResponseDeviation := time.Duration(0)

	// Start the timer
	startTime := time.Now()

	// Determine the output file using Viper
	outputFile := viper.GetString("tcp.ping.output-file")
	append := viper.GetBool("tcp.ping.append")

	// Get the output stream
	outputStream, err := utils.GetOutputStream(outputFile, append)
	if err != nil {
		return err
	}
	defer outputStream.Close()

	// Print start message (Initiate 3-way handshake with one.one.one.one (1.1.1.1) on port 443.)
	startMsg := fmt.Sprintf("Initiating 3-way handshakes with %s (%s) on port %d.\n", host, ip, port)

	// Print the compiled string to stdout
	fmt.Fprint(out, startMsg)

	// Print CSV header if --csv is set
	csvStartMsg := fmt.Sprintf("timestamp,host,ip,port,status,response_time_ms\n")

	// Print to file as well if --output-file is set
	if !viper.GetBool("tcp.ping.append") {
		if viper.IsSet("tcp.ping.output-file") && viper.GetBool("tcp.ping.csv") {
			fmt.Fprint(outputStream, csvStartMsg)
		} else if viper.IsSet("tcp.ping.output-file") {
			fmt.Fprint(outputStream, startMsg)
		}
	}

	// Start a goroutine that will print a message when a signal (Ctrl-C) is received
	go func() {
		sig := <-interrupt

		// Ctrl-C was pressed, print statistics and exit
		if sig == os.Interrupt {
			// Calculate mean deviation
			if packetsReceived > 1 {
				mdevResponseTime = totResponseDeviation / time.Duration(packetsReceived)
			}

			// Calculate total time
			totalTime := time.Since(startTime)
			totalTimeMs := totalTime.Round(time.Millisecond * 10)

			// Calculate min, avg, max and mdev response times
			avgResponseTimeMs := avgResponseTime.Round(time.Microsecond * 10)
			minResponseTimeMs := minResponseTime.Round(time.Microsecond * 10)
			maxResponseTimeMs := maxResponseTime.Round(time.Microsecond * 10)
			mdevResponseTimeMs := mdevResponseTime.Round(time.Microsecond * 10)

			// Calculate packet loss
			packetLoss := (packetsSent - packetsReceived) * 100 / packetsSent

			outStr := fmt.Sprintf("^C\n")
			outStr += fmt.Sprintf("--- %s ping statistics ---\n", host)
			outStr += fmt.Sprintf("%d packets transmitted, %d received, %d%% packet loss, time %s\n", packetsSent, packetsReceived, packetLoss, totalTimeMs)
			outStr += fmt.Sprintf("rtt min/avg/max/mdev = %s/%s/%s/%s\n", minResponseTimeMs, avgResponseTimeMs, maxResponseTimeMs, mdevResponseTimeMs)

			// Print the compiled string to stdout
			fmt.Fprint(out, outStr)

			// Print to file as well if --output-file is set and --csv is not set
			if viper.IsSet("tcp.ping.output-file") && !viper.GetBool("tcp.ping.csv") {
				fmt.Fprint(outputStream, outStr)
			}
			os.Exit(0)
		}
	}()

	// Set timeout duration for the TCP ping (default 2000 ms)
	timeoutMs := viper.GetDuration("tcp.ping.timeout") * time.Millisecond

	// Perform the TCP ping until user presses Ctrl-C
	for {
		// Send SYN packet and wait for SYN/ACK response
		packetsSent++

		// Send SYN packet and wait for SYN/ACK response
		responseTime, err := tcp.PingTCP(host, port, timeoutMs)

		// Check if the ping timed out
		if err != nil {
			// Get current time for timestamp
			currentTime := utils.GetTimestamp()

			// Format the CSV output string
			csvOutStr := fmt.Sprintf("%027s,%s,%s,%d,%s,%d\n", currentTime, host, ip, port, "offline", 0)

			// Print to file as well if --output-file is set
			if viper.IsSet("tcp.ping.output-file") && viper.GetBool("tcp.ping.csv") {
				fmt.Fprint(outputStream, csvOutStr)
			}

			if viper.GetBool("tcp.ping.verbose") {
				// Format the output string
				outStr := fmt.Sprintf("[%027s] Request timeout for %s: port=%d timeout=%s\n", currentTime, ip, port, timeoutMs)

				// Print the compiled string to stdout
				fmt.Fprint(out, outStr)

				// Print to file as well if --output-file is set
				if viper.IsSet("tcp.ping.output-file") && !viper.GetBool("tcp.ping.csv") {
					fmt.Fprint(outputStream, outStr)
				}
			} else {
				// Format the output string
				outStr := fmt.Sprintf("Request timeout for %s: port=%d timeout=%s\n", ip, port, timeoutMs)

				// Print the compiled string to stdout
				fmt.Fprint(out, outStr)

				// Print to file as well if --output-file is set
				if viper.IsSet("tcp.ping.output-file") && !viper.GetBool("tcp.ping.csv") {
					fmt.Fprint(outputStream, outStr)
				}
			}
			continue
		}

		// 3-way handshake completed, update packets received
		packetsReceived++

		// Update total response time
		totResponseTime += responseTime

		// Update min/max response times
		if packetsReceived == 1 {
			minResponseTime = responseTime
			maxResponseTime = responseTime
		} else {
			if responseTime < minResponseTime {
				minResponseTime = responseTime
			}
			if responseTime > maxResponseTime {
				maxResponseTime = responseTime
			}
		}

		// Update mean response time
		avgResponseTime = totResponseTime / time.Duration(packetsReceived)

		// Update mean deviation (mdev)
		// This is an average of how far each ping RTT is from the mean RTT. The higher mdev is, the more variable the RTT is (over time).
		stdResponseDeviation := float64(responseTime - avgResponseTime)
		stdResponseDeviation = math.Sqrt(math.Pow(stdResponseDeviation, 2))

		// Update total response deviation for later calculation of mdev
		totResponseDeviation += time.Duration(stdResponseDeviation)

		// Convert responseTime to float64
		responseTimeFloat := float64(responseTime) / float64(time.Millisecond)

		// Get current time for timestamp
		currentTime := utils.GetTimestamp()

		// Format the CSV output string
		csvOutStr := fmt.Sprintf("%s,%s,%s,%d,%s,%.4f\n", currentTime, host, ip, port, "online", responseTimeFloat)

		// Print to file as well if --output-file is set
		if viper.IsSet("tcp.ping.output-file") && viper.GetBool("tcp.ping.csv") {
			fmt.Fprint(outputStream, csvOutStr)
		}

		// Print response information (debug or normal output)
		if viper.GetBool("tcp.ping.verbose") {

			// Format the output string
			formatStr := "[%s] Received SYN/ACK from %s: port=%d tcp_seq=%d time=%-8s mrtt=%s\n"

			// Print to stdout
			fmt.Fprintf(out, formatStr, currentTime, ip, port, packetsSent, responseTime.Round(time.Microsecond*10), avgResponseTime.Round(time.Microsecond*10))

			// Print to file as well if --output-file is set
			if viper.IsSet("tcp.ping.output-file") && !viper.GetBool("tcp.ping.csv") {
				fmt.Fprintf(outputStream, formatStr, currentTime, ip, port, packetsSent, responseTime.Round(time.Microsecond*10), avgResponseTime.Round(time.Microsecond*10))
			}
		} else {
			// Format the output string
			formatStr := "Received SYN/ACK from %s: port=%d tcp_seq=%d time=%s\n"

			// Print to stdout
			fmt.Fprintf(out, formatStr, ip, port, packetsSent, responseTime.Round(time.Microsecond*10))

			// Print to file as well if --output-file is set
			if viper.IsSet("tcp.ping.output-file") && !viper.GetBool("tcp.ping.csv") {
				fmt.Fprintf(outputStream, formatStr, ip, port, packetsSent, responseTime.Round(time.Microsecond*10))
			}
		}

		// Check if the user specified a number of packets to send
		if count > 0 && packetsSent >= count {
			// Raise interrupt signal to stop the ping loop
			interrupt <- os.Interrupt
		}

		// Pause execution for the specified delay duration
		time.Sleep(delay)
	}
}

func init() {
	tcpCmd.AddCommand(pingCmd)

	// Enable the --timeout flag for the ping command
	pingCmd.Flags().IntP("timeout", "t", 2000, "time to wait for a response, in milliseconds")
	viper.BindPFlag("tcp.ping.timeout", pingCmd.Flags().Lookup("timeout"))

	// Enable the --delay flag for the ping command
	pingCmd.Flags().IntP("delay", "d", 1000, "delay between pings, in milliseconds")
	viper.BindPFlag("tcp.ping.delay", pingCmd.Flags().Lookup("delay"))

	// Enable the --count flag for the ping command
	pingCmd.Flags().IntP("count", "c", 0, "")
	viper.BindPFlag("tcp.ping.count", pingCmd.Flags().Lookup("count"))
	pingCmd.Flags().Lookup("count").Usage = "number of packets to send (default infinite)"

	// Enabled the --verbose flag for the ping command
	pingCmd.Flags().BoolP("verbose", "v", false, "show timestamps and mean round-trip time (mrtt)")
	viper.BindPFlag("tcp.ping.verbose", pingCmd.Flags().Lookup("verbose"))

	// Add flag for --output-file path
	pingCmd.PersistentFlags().StringP("output-file", "o", "", "write output to file")
	viper.BindPFlag("tcp.ping.output-file", pingCmd.PersistentFlags().Lookup("output-file"))

	// Set to the value of the --append flag if set
	pingCmd.PersistentFlags().BoolP("append", "a", false, "append when writing to file with --output-file")
	viper.BindPFlag("tcp.ping.append", pingCmd.PersistentFlags().Lookup("append"))

	// Set to the value of the --csv flag if set
	pingCmd.PersistentFlags().BoolP("csv", "C", false, "write output in CSV format")
	viper.BindPFlag("tcp.ping.csv", pingCmd.PersistentFlags().Lookup("csv"))

}
