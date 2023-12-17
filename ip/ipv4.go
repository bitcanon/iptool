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
package ip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// The IPv4 struct represents an IPv4 address as an IP address, a subnet mask
// and a network address. It also contains functions for calculating the
// broadcast address, the first and last usable host addresses, the number of
// usable hosts and the size of the network in number of IP addresses.
type IPv4 struct {
	IP   net.IP
	Mask net.IPMask
	Net  *net.IPNet
}

// Address is a function that returns the IP address in dotted-decimal notation
func (ip *IPv4) Address() string {
	return ip.IP.String()
}

// Netmask is a function that returns the netmask in dotted-decimal notation
func (ip *IPv4) Netmask() string {
	// Convert the hexadecimal string to an integer
	hexInt, err := strconv.ParseUint(ip.Mask.String(), 16, 32)
	if err != nil {
		return ""
	}

	// Convert the integer to dotted-decimal notation
	dottedDecimal := fmt.Sprintf("%d.%d.%d.%d",
		(hexInt>>24)&0xFF,
		(hexInt>>16)&0xFF,
		(hexInt>>8)&0xFF,
		hexInt&0xFF,
	)

	return dottedDecimal
}

// Wildcard is a function that returns the wildcard mask in dotted-decimal notation
func (ip *IPv4) Wildcard() string {
	// Convert the hexadecimal string to an integer
	hexInt, err := strconv.ParseUint(ip.Mask.String(), 16, 32)
	if err != nil {
		return ""
	}

	// Convert the integer to dotted-decimal notation
	dottedDecimal := fmt.Sprintf("%d.%d.%d.%d",

		(hexInt>>24)&0xFF^0xFF,
		(hexInt>>16)&0xFF^0xFF,
		(hexInt>>8)&0xFF^0xFF,
		hexInt&0xFF^0xFF,
	)

	return dottedDecimal
}

// Network is a function that returns the network address in the network
func (ip *IPv4) Network() string {
	return ip.Net.IP.String()
}

// PrefixLength is a function that returns the number of bits set in the netmask
func (ip *IPv4) PrefixLength() int {
	ones, _ := ip.Net.Mask.Size()
	return ones
}

// Broadcast is a function that returns the broadcast address in the network
func (ip *IPv4) Broadcast() string {
	// Convert the IP address to a 32-bit integer
	ipInt := ip.IP.To4()
	if ipInt == nil {
		return ""
	}
	ipInt32 := uint32(ipInt[0])<<24 | uint32(ipInt[1])<<16 | uint32(ipInt[2])<<8 | uint32(ipInt[3])

	// Convert the netmask to a 32-bit integer
	maskInt := ip.Mask
	maskInt32 := uint32(maskInt[0])<<24 | uint32(maskInt[1])<<16 | uint32(maskInt[2])<<8 | uint32(maskInt[3])

	// Calculate the broadcast address
	broadcastInt32 := ipInt32 | ^maskInt32

	// Convert the broadcast address to dotted-decimal notation
	broadcast := fmt.Sprintf("%d.%d.%d.%d",
		(broadcastInt32>>24)&0xFF,
		(broadcastInt32>>16)&0xFF,
		(broadcastInt32>>8)&0xFF,
		broadcastInt32&0xFF,
	)

	return broadcast
}

// FirstHost is a function that returns the first usable host address in the network
func (ip *IPv4) FirstHost() string {
	// Convert the IP address to a 32-bit integer
	ipInt := ip.IP.To4()
	if ipInt == nil {
		return ""
	}
	ipInt32 := uint32(ipInt[0])<<24 | uint32(ipInt[1])<<16 | uint32(ipInt[2])<<8 | uint32(ipInt[3])

	// Convert the netmask to a 32-bit integer
	maskInt := ip.Mask
	maskInt32 := uint32(maskInt[0])<<24 | uint32(maskInt[1])<<16 | uint32(maskInt[2])<<8 | uint32(maskInt[3])

	// Calculate the first host address
	firstHostInt32 := ipInt32 & maskInt32

	switch maskInt32 {
	// If maskInt32 is 0xFFFFFFFF, the network is a /32 network and the first host address is the same as the network address
	case 0xFFFFFFFF:
		firstHostInt32 = ipInt32
	// If maskInt32 is 0xFFFFFFFE, the network is a /31 network and the first host address is the same as the network address
	case 0xFFFFFFFE:
		firstHostInt32 = ipInt32
	// Else, the first host address is the network address + 1
	default:
		firstHostInt32 = firstHostInt32 + 1
	}

	// Convert the first host address to dotted-decimal notation
	firstHost := fmt.Sprintf("%d.%d.%d.%d",
		(firstHostInt32>>24)&0xFF,
		(firstHostInt32>>16)&0xFF,
		(firstHostInt32>>8)&0xFF,
		firstHostInt32&0xFF,
	)

	return firstHost
}

// LastHost is a function that returns the last usable host address in the network
func (ip *IPv4) LastHost() string {
	// Convert the IP address to a 32-bit integer
	ipInt := ip.IP.To4()
	if ipInt == nil {
		return ""
	}
	ipInt32 := uint32(ipInt[0])<<24 | uint32(ipInt[1])<<16 | uint32(ipInt[2])<<8 | uint32(ipInt[3])

	// Convert the netmask to a 32-bit integer
	maskInt := ip.Mask
	maskInt32 := uint32(maskInt[0])<<24 | uint32(maskInt[1])<<16 | uint32(maskInt[2])<<8 | uint32(maskInt[3])

	// Calculate the last host address
	lastHostInt32 := ipInt32 & maskInt32 //| ^maskInt32 - 1

	switch maskInt32 {
	// If maskInt32 is 0xFFFFFFFF, the network is a /32 network and the last host address is the same as the network and broadcast address
	case 0xFFFFFFFF:
		lastHostInt32 = ipInt32
	// If maskInt32 is 0xFFFFFFFE, the network is a /31 network and the last host address is the same as the broadcast address
	case 0xFFFFFFFE:
		lastHostInt32 = ipInt32 | ^maskInt32
	// Else, the last host address is the broadcast address - 1
	default:
		lastHostInt32 = lastHostInt32 | ^maskInt32 - 1
	}

	// Convert the last host address to dotted-decimal notation
	lastHost := fmt.Sprintf("%d.%d.%d.%d",
		(lastHostInt32>>24)&0xFF,
		(lastHostInt32>>16)&0xFF,
		(lastHostInt32>>8)&0xFF,
		lastHostInt32&0xFF,
	)

	return lastHost
}

// String is a function that returns the IP address and the prefix length in CIDR notation
func (ip *IPv4) String() string {
	return fmt.Sprintf("%s/%d", ip.IP.String(), ip.PrefixLength())
}

// UsableHosts is a function that returns the number of usable hosts in the network
func (ip *IPv4) UsableHosts() uint32 {
	// Convert the netmask to a 32-bit integer
	maskInt := ip.Mask
	maskInt32 := uint32(maskInt[0])<<24 | uint32(maskInt[1])<<16 | uint32(maskInt[2])<<8 | uint32(maskInt[3])

	// Get the number of bits set in the netmask
	ones, _ := ip.Net.Mask.Size()

	// In a /32 network, there are no usable hosts
	if ones == 32 {
		return 0
	}

	// In a /31 network, there are two usable hosts
	if ones == 31 {
		return 2
	}

	// Calculate the number of usable hosts
	usableHosts := ^maskInt32 - 1

	return usableHosts
}

// NetworkSize is a function that returns the size of the network in number of IP addresses
func (ip *IPv4) NetworkSize() uint32 {
	// Convert the netmask to a 32-bit integer
	maskInt := ip.Mask
	maskInt32 := uint32(maskInt[0])<<24 | uint32(maskInt[1])<<16 | uint32(maskInt[2])<<8 | uint32(maskInt[3])

	// Get the number of bits set in the netmask
	ones, _ := ip.Net.Mask.Size()

	// Calculate the network size
	networkSize := ^maskInt32 + 1

	// In a /0 network, the network size is 2^32 = 4294967296
	// But since we are using uint32, the maximum value is 4294967295
	if ones == 0 {
		return 4294967295
	}

	return networkSize
}

// NetmaskPrefixLength is a function that takes a netmask in dotted-decimal notation
// (e.g. 255.255.255.0) as input and returns the number of bits set in the netmask
func NetmaskPrefixLength(mask string) (int, error) {
	// Try to parse the netmask
	ip := net.ParseIP(mask)
	if ip == nil {
		return 0, ErrInvalidNetmask
	}

	// Count the number of bits set in the netmask
	ones, bits := net.IPMask(ip.To4()).Size()

	// Make sure that the netmask is a valid IPv4 netmask
	if bits != 32 {
		return 0, ErrInvalidNetmask
	}

	// Return the number of bits set in the netmask
	return ones, nil
}

// ParseIPv4 is a function that takes a string as input and returns an IPv4 address
// and a subnet mask as output.
// The input string can be in the following formats:
// - "X.X.X.X/Y"
// - "X.X.X.X Y"
// - "X.X.X.X"
func ParseIPv4(s string) (*IPv4, error) {
	// Try to split the input string into an IP address and a netmask
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '/' || r == ' '
	})

	// If the input string contains two parts, check if the second part is a netmask
	// in dotted-decimal notation (255.255.255.0) or CIDR notation (24)
	if len(parts) == 2 {
		// If the netmask is in dotted-decimal notation, convert it to CIDR notation
		if IsIPv4(parts[1]) {
			ones, err := NetmaskPrefixLength(parts[1])
			if err != nil {
				return nil, err
			}
			parts[1] = strconv.Itoa(ones)
		}
	} else if len(parts) == 1 {
		// If the input string does not contain a netmask or prefix length,
		// assume that the netmask is 24 bits
		parts = append(parts, "24")
	} else {
		return nil, fmt.Errorf("invalid IP address: %s", s)
	}

	// Reassemble the input string
	s = strings.Join(parts, "/")

	// Parse the input string
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &IPv4{IP: ip, Mask: ipnet.Mask, Net: ipnet}, nil
}
