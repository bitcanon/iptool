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
	"os"
	"strings"

	"github.com/bitcanon/iptool/debug"
	"github.com/bitcanon/iptool/ip"
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
	// networks := viper.GetInt("subnet.split.networks")

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

	// Print the subnets
	// Start with the header (Prefix, Network, Broadcast, First, Last, Hosts)
	fmt.Fprintf(out, fmtString, "Prefix", "Network", "First", "Last", "Broadcast", "Hosts")
	fmt.Println(dashLine)
	for _, prefix := range prefixList {
		pfx := prefix.String()
		network := prefix.Network()
		broadcast := prefix.Broadcast()
		first := prefix.FirstHost()
		last := prefix.LastHost()
		hosts := prefix.UsableHosts()

		fmt.Fprintf(out, fmtString, pfx, network, first, last, broadcast, fmt.Sprint(hosts))
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
	subnetSplitCmd.Flags().IntP("bits", "b", 30, "subnet size in bits for network division")
	viper.BindPFlag("subnet.split.bits", subnetSplitCmd.Flags().Lookup("bits"))

	// Define the flag for specifying the number of subnets to split the network into
	subnetSplitCmd.Flags().IntP("networks", "n", 0, "number of subnets to divide the network into")
	viper.BindPFlag("subnet.split.networks", subnetSplitCmd.Flags().Lookup("networks"))
}
