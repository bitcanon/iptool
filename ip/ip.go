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
	"errors"
	"net"
	"strings"
)

var ErrInvalidNetmask = errors.New("invalid netmask")

// Function that takes a string as input and returns an IP address
// and a subnet mask as output.
func ParseIP(s string) (net.IP, *net.IPNet, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, nil, err
	}
	return ip, ipnet, nil
}

// Function that checks if an IP address is an IPv4 address
func IsIPv4(s string) bool {
	ipv4 := net.ParseIP(s)
	if ipv4 != nil && strings.Contains(s, ".") {
		return true
	}
	return false
}

// IsIPv6 is a function that checks if an IP address is an IPv6 address
func IsIPv6(s string) bool {
	ipv6 := net.ParseIP(s)
	if ipv6 != nil && strings.Contains(s, ":") {
		return true
	}
	return false
}
