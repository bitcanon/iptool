/*
Copyright © 2023 Mikael Schultz <bitcanon@proton.me>

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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping <destination> [port]",
	Short: "Send a stream of TCP pings to a host",
	Long: `Send a stream of TCP pings to a host.

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

	// Resolve the IP address of the destination
	ip, err := ip.ResolveIP(host)
	if err != nil {
		return err
	}

	// Print start message (Initiate 3-way handshake with one.one.one.one (1.1.1.1) on port 443.)
	fmt.Fprintf(out, "Initiating 3-way handshakes with %s (%s) on port %d.\n", host, ip, port)

	// Create a channel to receive interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Initialize variables for ping statistics
	packetsSent := 0
	packetsReceived := 0
	packetLoss := 0
	minResponseTime := time.Duration(0)
	maxResponseTime := time.Duration(0)
	avgResponseTime := time.Duration(0)
	mdevResponseTime := time.Duration(0)
	totResponseTime := time.Duration(0)
	totResponseDeviation := time.Duration(0)
	tcpSeq := 0

	// Start the timer
	startTime := time.Now()

	// Start a goroutine that will print a message when a signal is received
	go func() {
		sig := <-interrupt

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
			packetLoss = (packetsSent - packetsReceived) * 100 / packetsSent

			fmt.Fprintln(out, "^C")
			fmt.Fprintf(out, "--- %s ping statistics ---\n", host)
			fmt.Fprintf(out, "%d packets transmitted, %d received, %d%% packet loss, time %s\n", packetsSent, packetsReceived, packetLoss, totalTimeMs)
			fmt.Fprintf(out, "rtt min/avg/max/mdev = %s/%s/%s/%s\n", minResponseTimeMs, avgResponseTimeMs, maxResponseTimeMs, mdevResponseTimeMs)
			os.Exit(0)
		}
	}()

	// Set timeout duration for the TCP ping (default 2000 ms)
	timeoutMs := viper.GetDuration("tcp.ping.timeout") * time.Millisecond

	// Perform the TCP ping until user presses Ctrl-C
	for {
		// Send SYN packet and wait for SYN/ACK response
		packetsSent++

		responseTime, err := tcp.PingTCP(host, port, timeoutMs)
		if err != nil {
			fmt.Fprintf(out, "Request timeout for %s: port=%d timeout=%s\n", ip, port, timeoutMs)
			continue
		}
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

		// Print response information (debug or normal output)
		if viper.GetBool("tcp.ping.verbose") {
			currentTime := time.Now().Format("2006-01-02 15:04:05.999999999")

			fmt.Fprintf(out, "[%-27s] Received SYN/ACK from %s: port=%d tcp_seq=%d time=%-8s mrtt=%s\n", currentTime, ip, port, packetsSent, responseTime.Round(time.Microsecond*10), avgResponseTime.Round(time.Microsecond*10))
		} else {
			fmt.Fprintf(out, "Received SYN/ACK from %s: port=%d tcp_seq=%d time=%s\n", ip, port, packetsSent, responseTime.Round(time.Microsecond*10))
		}

		// Update TCP sequence number
		tcpSeq++

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
}
