/*
Copyright © 2024 Mikael Schultz <bitcanon@proton.me>

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
	"strings"

	"github.com/bitcanon/iptool/debug"
	"github.com/bitcanon/iptool/ip"
	"github.com/bitcanon/iptool/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// subnetSplitCmd represents the subnetSplit command
var subnetSplitCmd = &cobra.Command{
	Use:   "split <subnet>",
	Short: "Splits a given subnet into smaller subnets",
	Long: `Splits a given subnet into smaller subnets based on the specified size or number of subnets.

Examples:
  iptool subnet split 10.0.0.0/24 --bits 30
  iptool subnet split 10.0.0.0 255.255.255.0 --networks 4`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no arguments are provided, print a short help text
		if len(args) == 0 {
			cmd.Help()
			return nil
		}
		input := args[0]

		return subnetSplitAction(os.Stdout, input)
	},
}

// subnetSplitAction is the action function for the subnetSplit command
func subnetSplitAction(out io.Writer, s string) error {
	// Parse the input string as an IP address
	network, err := ip.ParseIPv4(s)
	if err != nil {
		return err
	}

	// Parse the network count and bits from the configuration
	bits := viper.GetInt("subnet.split.bits")
	networks := viper.GetInt("subnet.split.networks")

	// If both bits and networks are specified, return an error
	if bits > 0 && networks > 0 {
		return fmt.Errorf("both --bits and --networks cannot be specified at the same time, see --help for more information")
	}

	// If neither bits nor networks are specified, return an error
	if bits == 0 && networks == 0 {
		return fmt.Errorf("either --bits or --networks must be specified, see --help for more information")
	}

	// If the number of networks is specified, calculate the number of bits required
	if networks > 0 {
		// Calculate the number of networks closest to a power of two (2, 4, 8, 16, 32, 64, 128, 256, ...)
		subnets := utils.ClosestLargerPowerOfTwo(networks)

		// Calculate the number of addresses in each subnet
		subnetSize := network.NetworkSize() / subnets

		// Calculate the number of host bits in each subnet
		hostBits := int(math.Log2(float64(subnetSize)))

		// Calculate the number of network bits
		bits = 32 - hostBits
	}

	// Split the subnet into smaller subnets
	prefixList, err := network.Split(bits)
	if err != nil {
		return err
	}

	// Find the length of the longest broadcast address (for padding)
	// This is used to align Prefix, Network, Broadcast, First, Last, Hosts
	maxLength := 0
	for _, prefix := range prefixList {
		broadcast := prefix.Broadcast()
		if len(broadcast) > maxLength {
			maxLength = len(broadcast)
		}
	}
	maxLength += 1

	// Format string for padding
	fmtString := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%s\n", maxLength+3, maxLength, maxLength, maxLength, maxLength)

	// Calculate the total length of the output
	columns := 5
	spacesBetweenColumns := 2 * columns
	totalLength := (maxLength * columns) + spacesBetweenColumns + 3

	// Create a string of dashes of the total length
	dashLine := strings.Repeat("-", totalLength)

	// Determine the output file using Viper
	outputFile := viper.GetString("subnet.split.output-file")

	// Get the output stream
	outputStream, err := utils.GetOutputStream(outputFile, false)
	if err != nil {
		return err
	}
	defer outputStream.Close()

	// Print the subnets
	// Start with the header (Prefix, Network, Broadcast, First, Last, Hosts)
	if viper.GetBool("subnet.split.csv") {
		fmt.Fprintf(outputStream, "prefix,network,first,last,broadcast,hosts\n")
	} else {
		fmt.Fprintf(outputStream, fmtString, "Prefix", "Network", "First", "Last", "Broadcast", "Hosts")
		fmt.Fprintf(outputStream, dashLine+"\n")
	}
	for _, prefix := range prefixList {
		pfx := prefix.String()
		network := prefix.Network()
		broadcast := prefix.Broadcast()
		first := prefix.FirstHost()
		last := prefix.LastHost()
		hosts := prefix.UsableHosts()

		if viper.GetBool("subnet.split.csv") {
			fmt.Fprintf(outputStream, "%s,%s,%s,%s,%s,%s\n", pfx, network, first, last, broadcast, fmt.Sprint(hosts))
		} else {
			fmt.Fprintf(outputStream, fmtString, pfx, network, first, last, broadcast, fmt.Sprint(hosts))
		}
	}

	// Print the configuration debug if the --debug flag is set
	if viper.GetBool("debug") {
		debug.PrintConfigDebug()
	}

	return nil
}

func init() {
	subnetCmd.AddCommand(subnetSplitCmd)

	// Define the flag for specifying the size of the subnets
	subnetSplitCmd.Flags().IntP("bits", "b", 0, "subnet size in bits for network division")
	viper.BindPFlag("subnet.split.bits", subnetSplitCmd.Flags().Lookup("bits"))

	// Define the flag for specifying the number of subnets to split the network into
	subnetSplitCmd.Flags().IntP("networks", "n", 0, "number of subnets to divide the network into")
	viper.BindPFlag("subnet.split.networks", subnetSplitCmd.Flags().Lookup("networks"))

	// Define the flag for allowing the user to output in CSV format
	subnetSplitCmd.Flags().BoolP("csv", "c", false, "output in CSV format")
	viper.BindPFlag("subnet.split.csv", subnetSplitCmd.Flags().Lookup("csv"))

	// Define the flag for allowing the user to output to a file
	subnetSplitCmd.Flags().StringP("output-file", "o", "", "write output to file")
	viper.BindPFlag("subnet.split.output-file", subnetSplitCmd.Flags().Lookup("output-file"))
}
