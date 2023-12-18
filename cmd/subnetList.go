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
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bitcanon/iptool/debug"
	"github.com/bitcanon/iptool/ip"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// subnetListCmd represents the subnetList command
var subnetListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a comprehensive IPv4 subnet mask list",
	Long: `Display a comprehensive IPv4 subnet mask list, ranging from 0 to 32 bits.

Filter the list by specifying one or more prefix lengths (integers
between 0 and 32) as an argument, separated by commas.

Examples:
  iptool subnet list
  iptool subnet list -p 8,16,24
`,
	Aliases:      []string{"ls"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// No arguments allowed
		if len(args) > 0 {
			return fmt.Errorf("invalid argument(s): %s", strings.Join(args, " "))
		}

		input := strings.Join(args, " ")
		return subnetListAction(os.Stdout, input)
	},
}

// subnetListAction prints a list of IPv4 subnets
func subnetListAction(out io.Writer, s string) error {
	// Print the header for the table
	fmt.Fprintf(out, "CIDR  Subnet Mask      Addresses   Wildcard Mask\n")
	fmt.Fprintf(out, "--------------------------------------------------\n")

	// Get the prefix lengths from the viper configuration
	prefixList := viper.GetIntSlice("subnet.list.prefix-lengths")

	// If prefixList is empty, add all prefix lengths (0-32)
	if len(prefixList) == 0 {
		for i := 32; i >= 0; i-- {
			prefixList = append(prefixList, i)
		}
	}

	// Loop through all subnets
	for _, i := range prefixList {
		// Print information about the subnet
		s = fmt.Sprintf("0.0.0.0/%d", i)
		subnet, err := ip.ParseIPv4(s)
		if err != nil {
			return err
		}

		// Print information about the subnet
		fmt.Fprintf(out, "%4s  %-16s %-11d %-10s\n", "/"+strconv.Itoa(subnet.PrefixLength()), subnet.Netmask(), subnet.NetworkSize(), subnet.Wildcard())
	}

	// Print the configuration debug if the --debug flag is set
	if viper.GetBool("debug") {
		debug.PrintConfigDebug()
	}

	return nil
}

// init registers the command and flags
func init() {
	subnetCmd.AddCommand(subnetListCmd)

	// Define the flag for prefix lengths
	subnetListCmd.Flags().IntSliceP("prefix-lengths", "p", []int{}, "a list of prefix lengths (0-32)")
	viper.BindPFlag("subnet.list.prefix-lengths", subnetListCmd.Flags().Lookup("prefix-lengths"))

	// Validate the prefix lengths
	subnetListCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		for _, length := range viper.GetIntSlice("subnet.list.prefix-lengths") {
			if length < 0 || length > 32 {
				message := fmt.Sprintf("invalid prefix length: %d (must be between 0 and 32)", length)
				return errors.New(message)
			}
		}
		return nil
	}
}
