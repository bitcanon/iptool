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
	"github.com/spf13/viper"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect <ip address>",
	Short: "Take a closer look at an IP address",
	Long: `Inspect an IP address in any format and print detailed information about
the address. If no subnet mask is specified, a subnet mask of 24 bits is assumed.

Examples:
  iptool inspect 10.0.0.1
  iptool inspect 10.0.0.1/24
  iptool inspect 10.0.0.1 255.255.255.0
  iptool inspect 0xc0800d25
  iptool inspect c0800d25/22
  iptool inspect c0800d25 fffffe00`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no arguments are provided, print a short help text
		if len(args) == 0 {
			cmd.Help()
			return nil
		}
		input := strings.Join(args, " ")

		return inspectAction(os.Stdout, input)
	},
}

const simpleTemplate = `Address Details:
 IPv4 address       : {{.HostAddress}}
 Network mask       : {{.NetworkMask}}

Netmask Details:
 Network mask       : {{.NetworkMask}}
 Network bits       : {{.NetworkMaskBits}}
 Wildcard mask      : {{.WildcardMask}}

Network Details:
 CIDR notation      : {{.NetworkDetails}} ({{.NetworkSize}} addresses)
 Network address    : {{.NetworkAddress}}
 Broadcast address  : {{.BroadcastAddress}}
 Usable hosts       : {{.FirstHost}} - {{.LastHost}} ({{.UsableHosts}} hosts)
`

const advancedTemplate = `Address Details:
 IPv4 address       : {{.HostAddress}}
 Network mask       : {{.NetworkMask}}

Netmask Details:
 Network mask       : {{.NetworkMask}}
 Network bits       : {{.NetworkMaskBits}}
 Wildcard mask      : {{.WildcardMask}}

Network Details:
 CIDR notation      : {{.NetworkDetails}} ({{.NetworkSize}} addresses)
 Network address    : {{.NetworkAddress}}
 Broadcast address  : {{.BroadcastAddress}}
 Usable hosts       : {{.FirstHost}} - {{.LastHost}} ({{.UsableHosts}} hosts)

Binary Notation:
 IPv4 address       : {{.HostAddressBinary}} ({{.HostAddress}})
 Network mask       : {{.NetworkMaskBinary}} ({{.NetworkMask}})
 Network address    : {{.NetworkAddressBinary}} ({{.NetworkAddress}})
 Broadcast address  : {{.BroadcastAddressBinary}} ({{.BroadcastAddress}})
 Wildcard mask      : {{.WildcardMaskBinary}} ({{.WildcardMask}})

Hexadecimal Notation:
 IPv4 address       : {{.HostAddressHex}} ({{.HostAddress}})
 Network mask       : {{.NetworkMaskHex}} ({{.NetworkMask}})
 Network address    : {{.NetworkAddressHex}} ({{.NetworkAddress}})
 Broadcast address  : {{.BroadcastAddressHex}} ({{.BroadcastAddress}})
 Wildcard mask      : {{.WildcardMaskHex}} ({{.WildcardMask}})

Decimal Notation:
 IPv4 address       : {{printf "%10s" .HostAddressDecimal}} ({{.HostAddress}})
 Network mask       : {{printf "%10s" .NetworkMaskDecimal}} ({{.NetworkMask}})
 Network address    : {{printf "%10s" .NetworkAddressDecimal}} ({{.NetworkAddress}})
 Broadcast address  : {{printf "%10s" .BroadcastAddressDecimal}} ({{.BroadcastAddress}})
 Wildcard mask      : {{printf "%10s" .WildcardMaskDecimal}} ({{.WildcardMask}})
`

func inspectAction(out io.Writer, s string) error {
	if strings.Contains(s, ":") {
		// If there is a colon in the input string, assume it is an IPv6 address
		return fmt.Errorf("support for IPv6 addresses is not implemented yet")
	} else {
		// Otherwise, assume it is an IPv4 address (either in hexadecimal or dotted decimal notation)
		ipv4, err := ip.ParseIPv4(s)
		if err != nil {
			return err
		}

		// Create a data structure with the values to fill in the template placeholders
		data := struct {
			NetworkMask             string
			NetworkMaskBinary       string
			NetworkMaskHex          string
			NetworkMaskDecimal      string
			NetworkDetails          string
			HostAddress             string
			HostAddressBinary       string
			HostAddressHex          string
			HostAddressDecimal      string
			NetworkAddress          string
			NetworkAddressBinary    string
			NetworkAddressHex       string
			NetworkAddressDecimal   string
			BroadcastAddress        string
			BroadcastAddressBinary  string
			BroadcastAddressHex     string
			BroadcastAddressDecimal string
			UsableHosts             string
			FirstHost               string
			LastHost                string
			NetworkSize             string
			NetworkMaskBits         string
			WildcardMask            string
			WildcardMaskBinary      string
			WildcardMaskHex         string
			WildcardMaskDecimal     string
		}{
			NetworkMask:             ipv4.Netmask(),
			NetworkMaskBinary:       ip.IPv4ToBinary(ipv4.Netmask()),
			NetworkMaskHex:          ip.IPv4ToHex(ipv4.Netmask()),
			NetworkMaskDecimal:      ip.IPv4ToDecimal(ipv4.Netmask()),
			NetworkDetails:          fmt.Sprintf("%s/%d", ipv4.Network(), ipv4.PrefixLength()),
			HostAddress:             ipv4.Address(),
			HostAddressBinary:       ip.IPv4ToBinary(ipv4.Address()),
			HostAddressHex:          ip.IPv4ToHex(ipv4.Address()),
			HostAddressDecimal:      ip.IPv4ToDecimal(ipv4.Address()),
			NetworkAddress:          ipv4.Network(),
			NetworkAddressBinary:    ip.IPv4ToBinary(ipv4.Network()),
			NetworkAddressHex:       ip.IPv4ToHex(ipv4.Network()),
			NetworkAddressDecimal:   ip.IPv4ToDecimal(ipv4.Network()),
			BroadcastAddress:        ipv4.Broadcast(),
			BroadcastAddressBinary:  ip.IPv4ToBinary(ipv4.Broadcast()),
			BroadcastAddressHex:     ip.IPv4ToHex(ipv4.Broadcast()),
			BroadcastAddressDecimal: ip.IPv4ToDecimal(ipv4.Broadcast()),
			UsableHosts:             fmt.Sprintf("%d", ipv4.UsableHosts()),
			FirstHost:               ipv4.FirstHost(),
			LastHost:                ipv4.LastHost(),
			NetworkSize:             fmt.Sprintf("%d", ipv4.NetworkSize()),
			NetworkMaskBits:         fmt.Sprintf("%d", ipv4.PrefixLength()),
			WildcardMask:            ipv4.Wildcard(),
			WildcardMaskBinary:      ip.IPv4ToBinary(ipv4.Wildcard()),
			WildcardMaskHex:         ip.IPv4ToHex(ipv4.Wildcard()),
			WildcardMaskDecimal:     ip.IPv4ToDecimal(ipv4.Wildcard()),
		}

		// If the --detailed flag is set, use the advanced template
		selectedTemplate := simpleTemplate
		if viper.GetBool("inspect.detailed") {
			selectedTemplate = advancedTemplate
		}

		// Create a new template and parse the template text
		tmpl := template.Must(template.New("networkDetails").Parse(selectedTemplate))

		// Execute the template with the data and write the result to an output
		err = tmpl.Execute(os.Stdout, data)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
	}

	return nil
}

func init() {
	// Register the inspect command with the root command
	rootCmd.AddCommand(inspectCmd)

	// Enable the --detailed flag for the inspect command
	inspectCmd.Flags().BoolP("detailed", "d", false, "display comprehensive IP address information")
	viper.BindPFlag("inspect.detailed", inspectCmd.Flags().Lookup("detailed"))
}
