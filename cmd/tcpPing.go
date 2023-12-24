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
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"strconv"
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
  tcp ping 1.0.0.1
  tcp ping 1.0.0.1 443`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// // Parse the arguments (destination and optional port)
		// if len(args) < 1 || len(args) > 2 {
		// 	return errors.New("invalid number of arguments")
		// }

		// Parse the destination
		destination := args[0]

		// Parse the port
		port := 443
		if len(args) == 2 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			port = p
		}

		return tcpPingAction(os.Stdout, destination, port)
	},
}

func tcpPingAction(out io.Writer, host string, port int) error {
	// Define the delay duration
	var delay time.Duration = 1 * time.Second

	if viper.IsSet("tcp.ping.delay") {
		delay = viper.GetDuration("tcp.ping.delay")
	}

	// Print delay information (Delay is 1 second)
	fmt.Fprintf(out, "Delay is %s\n", delay)

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

	// Set timeout duration
	timeoutMs := viper.GetDuration("tcp.ping.timeout") * time.Millisecond

	// Perform the TCP ping until user presses Ctrl-C
	for {
		// Send SYN packet and wait for SYN/ACK response
		packetsSent++

		responseTime, err := tcp.PingTCP(host, port, 10, timeoutMs)
		if err != nil {
			fmt.Fprintf(out, "Request timeout for %s: port=%d ttl=12 timeout=%s\n", ip, port, timeoutMs)
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
		stdResponseDeviationMs := time.Duration(stdResponseDeviation)

		// Update total response deviation for later calculation of mdev
		totResponseDeviation += time.Duration(stdResponseDeviation)

		// Print response information (debug or normal output)
		if viper.GetBool("debug") {
			fmt.Fprintf(out, "Received SYN/ACK from %s: port=%d tcp_seq=%d ttl=54 time=%s mrtt=%s dev=%s\n", ip, port, packetsSent, responseTime.Round(time.Microsecond*10), avgResponseTime, stdResponseDeviationMs)
		} else {
			fmt.Fprintf(out, "Received SYN/ACK from %s: port=%d tcp_seq=%d ttl=54 time=%s\n", ip, port, packetsSent, responseTime.Round(time.Microsecond*10))
		}

		// Pause execution for the specified delay duration
		time.Sleep(delay)

		// Update TCP sequence number
		tcpSeq++

		// [TODO] totalTime should show the total time since the first packet was sent
		// [TODO] When pinging an address that is not responding, print "Request timeout for 65.9.55.64: port=443 ttl=12 time=1.05sec"
		// [TODO] Calculate min, avg, max and mdev response times
		// [TODO] Implement tcp.ping.delay flag
		// [TODO] Implement tcp.ping.count flag
		// [TODO] Implement tcp.ping.interval flag
		// [TODO] Implement tcp.ping.timeout flag
		// [TODO] Implement tcp.ping.ttl flag

	}
}

func init() {
	tcpCmd.AddCommand(pingCmd)

	// Enable the --detailed flag for the inspect command
	pingCmd.Flags().IntP("timeout", "t", 1000, "time in milliseconds to wait for a response")
	viper.BindPFlag("tcp.ping.timeout", pingCmd.Flags().Lookup("timeout"))
}
