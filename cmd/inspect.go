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
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/bitcanon/iptool/ip"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect <ip address>",
	Short: "Inspect an IP address in any format.",
	Long: `Inspect an IP address in any format and print detailed information about
the address. If no subnet mask is specified, a subnet mask of 24 bits is assumed.

Examples:
  iptool inspect 10.0.0.1
  iptool inspect 10.0.0.1/24
  iptool inspect 10.0.0.1 255.255.255.0`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.Join(args, " ")
		return inspectAction(os.Stdout, input)
	},
}

func inspectAction(out io.Writer, s string) error {
	templateText := `Network details: {{.NetworkDetails}} ({{.NetworkSize}} addresses)

Address Details
---------------
IPv4 address       : {{.HostAddress}} {{.NetworkMask}}
Network address    : {{.NetworkAddress}}
Broadcast address  : {{.BroadcastAddress}}
Usable hosts	   : {{.FirstHost}} - {{.LastHost}} ({{.UsableHosts}} hosts)

Network Mask
------------
Network mask         : {{.NetworkMask}}
Network mask (bits)  : {{.NetworkMaskBits}}
`

	// Check if the input is an IPv4 or IPv6 address
	if strings.Contains(s, ".") {
		// Parse the IP address and subnet mask
		ipv4, err := ip.ParseIPv4(s)
		if err != nil {
			return err
		}

		// Create a data structure with the values to fill in the template placeholders
		data := struct {
			NetworkMask      string
			NetworkDetails   string
			HostAddress      string
			NetworkAddress   string
			BroadcastAddress string
			UsableHosts      string
			FirstHost        string
			LastHost         string
			NetworkSize      string
			NetworkMaskBits  string
		}{
			NetworkMask:      ipv4.Netmask(),
			NetworkDetails:   fmt.Sprintf("%s/%d", ipv4.Network(), ipv4.PrefixLength()),
			HostAddress:      ipv4.Address(),
			NetworkAddress:   ipv4.Network(),
			BroadcastAddress: ipv4.Broadcast(),
			UsableHosts:      fmt.Sprintf("%d", ipv4.UsableHosts()),
			FirstHost:        ipv4.FirstHost(),
			LastHost:         ipv4.LastHost(),
			NetworkSize:      fmt.Sprintf("%d", ipv4.NetworkSize()),
			NetworkMaskBits:  fmt.Sprintf("%d", ipv4.PrefixLength()),
		}

		// Create a new template and parse the template text
		tmpl := template.Must(template.New("networkDetails").Parse(templateText))

		// Execute the template with the data and write the result to an output
		err = tmpl.Execute(os.Stdout, data)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
	} else {
		// In the case of an IPv6 address, we need to parse the IP address and prefix length
		// To be implemented in a future version
		return fmt.Errorf("invalid IP address: %s", s)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
