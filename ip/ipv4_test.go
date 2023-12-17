package ip_test

import (
	"testing"

	"github.com/bitcanon/iptool/ip"
)

// TestParseIPv4 is a function that tests the ParseIPv4 function.
func TestParseIPv4(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name             string
		input            string
		expectedIP       string
		expectedMask     string
		expectedMaskBits int
		expectedNet      string
		expectedCIDR     string
	}{
		{
			name:             "IPv4AddresInCIDRNotation",
			input:            "1.2.3.4/24",
			expectedIP:       "1.2.3.4",
			expectedMask:     "255.255.255.0",
			expectedMaskBits: 24,
			expectedNet:      "1.2.3.0",
			expectedCIDR:     "1.2.3.4/24",
		},
		{
			name:             "IPv4AddresInCIDRNotation",
			input:            "1.2.3.4/22",
			expectedIP:       "1.2.3.4",
			expectedMask:     "255.255.252.0",
			expectedMaskBits: 22,
			expectedNet:      "1.2.0.0",
			expectedCIDR:     "1.2.3.4/22",
		},
		{
			name:             "IPv4AddressWithoutNetmask",
			input:            "1.2.3.4",
			expectedIP:       "1.2.3.4",
			expectedMask:     "255.255.255.0",
			expectedMaskBits: 24,
			expectedNet:      "1.2.3.0",
			expectedCIDR:     "1.2.3.4/24",
		},
		{
			name:             "IPv4AddressWithNetmask",
			input:            "1.2.3.4 255.255.255.0",
			expectedIP:       "1.2.3.4",
			expectedMask:     "255.255.255.0",
			expectedMaskBits: 24,
			expectedNet:      "1.2.3.0",
			expectedCIDR:     "1.2.3.4/24",
		},
		{
			name:             "IPv4DefaultRoute",
			input:            "0.0.0.0 0.0.0.0",
			expectedIP:       "0.0.0.0",
			expectedMask:     "0.0.0.0",
			expectedMaskBits: 0,
			expectedNet:      "0.0.0.0",
			expectedCIDR:     "0.0.0.0/0",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ipv4, err := ip.ParseIPv4(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ipv4.Address() != tc.expectedIP {
				t.Errorf("expected IP address %q, got %q", tc.expectedIP, ipv4.Address())
			}
			if ipv4.Netmask() != tc.expectedMask {
				t.Errorf("expected netmask %q, got %q", tc.expectedMask, ipv4.Netmask())
			}
			if ipv4.PrefixLength() != tc.expectedMaskBits {
				t.Errorf("expected netmask prefix length %d, got %d", tc.expectedMaskBits, ipv4.PrefixLength())
			}
			if ipv4.String() != tc.expectedCIDR {
				t.Errorf("expected CIDR %q, got %q", tc.expectedCIDR, ipv4.String())
			}
			if ipv4.Network() != tc.expectedNet {
				t.Errorf("expected network %q, got %q", tc.expectedNet, ipv4.Network())
			}
		})
	}
}

func TestIPv4Broadcast(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Slash32", input: "10.0.0.1/32", expected: "10.0.0.1"},
		{name: "Slash31", input: "10.0.0.1/31", expected: "10.0.0.1"},
		{name: "Slash30", input: "10.0.0.1/30", expected: "10.0.0.3"},
		{name: "Slash24", input: "10.0.0.1/24", expected: "10.0.0.255"},
		{name: "Slash22", input: "10.0.0.1/22", expected: "10.0.3.255"},
		{name: "Slash30Net", input: "10.0.0.0/30", expected: "10.0.0.3"},
		{name: "Slash30First", input: "10.0.0.1/30", expected: "10.0.0.3"},
		{name: "Slash30Last", input: "10.0.0.2/30", expected: "10.0.0.3"},
		{name: "Slash30Bcast", input: "10.0.0.3/30", expected: "10.0.0.3"},
		{name: "Slash1", input: "10.0.0.1/1", expected: "127.255.255.255"},
		{name: "Slash0", input: "10.0.0.1/0", expected: "255.255.255.255"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ipv4, err := ip.ParseIPv4(tc.input)

			// Check for unexpected error
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ipv4.Broadcast() != tc.expected {
				t.Errorf("expected broadcast address %q, got %q", tc.expected, ipv4.Broadcast())
			}
		})
	}
}

func TestIPv4FirstHost(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Slash32", input: "10.0.0.1/32", expected: "10.0.0.1"},
		{name: "Slash31", input: "10.0.0.1/31", expected: "10.0.0.1"},
		{name: "Slash30", input: "10.0.0.1/30", expected: "10.0.0.1"},
		{name: "Slash24", input: "10.0.0.1/24", expected: "10.0.0.1"},
		{name: "Slash22", input: "10.0.0.1/22", expected: "10.0.0.1"},
		{name: "Slash30Net", input: "10.0.0.0/30", expected: "10.0.0.1"},
		{name: "Slash30First", input: "10.0.0.1/30", expected: "10.0.0.1"},
		{name: "Slash30Last", input: "10.0.0.2/30", expected: "10.0.0.1"},
		{name: "Slash30Bcast", input: "10.0.0.3/30", expected: "10.0.0.1"},
		{name: "Slash1", input: "10.0.0.1/1", expected: "0.0.0.1"},
		{name: "Slash0", input: "10.0.0.1/0", expected: "0.0.0.1"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ipv4, err := ip.ParseIPv4(tc.input)

			// Check for unexpected error
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ipv4.FirstHost() != tc.expected {
				t.Errorf("expected first host %q, got %q", tc.expected, ipv4.FirstHost())
			}
		})
	}
}

func TestIPv4LastHost(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Slash32", input: "10.0.0.1/32", expected: "10.0.0.1"},
		{name: "Slash31", input: "10.0.0.1/31", expected: "10.0.0.1"},
		{name: "Slash30", input: "10.0.0.1/30", expected: "10.0.0.2"},
		{name: "Slash24", input: "10.0.0.1/24", expected: "10.0.0.254"},
		{name: "Slash22", input: "10.0.0.1/22", expected: "10.0.3.254"},
		{name: "Slash30Net", input: "10.0.0.0/30", expected: "10.0.0.2"},
		{name: "Slash30First", input: "10.0.0.1/30", expected: "10.0.0.2"},
		{name: "Slash30Last", input: "10.0.0.2/30", expected: "10.0.0.2"},
		{name: "Slash30Bcast", input: "10.0.0.3/30", expected: "10.0.0.2"},
		{name: "Slash1", input: "10.0.0.1/1", expected: "127.255.255.254"},
		{name: "Slash0", input: "10.0.0.1/0", expected: "255.255.255.254"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ipv4, err := ip.ParseIPv4(tc.input)

			// Check for unexpected error
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ipv4.LastHost() != tc.expected {
				t.Errorf("expected last host %q, got %q", tc.expected, ipv4.LastHost())
			}
		})
	}
}

func TestIPv4String(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Slash32", input: "10.0.0.1/32", expected: "10.0.0.1/32"},
		{name: "Slash31", input: "10.0.0.1/31", expected: "10.0.0.1/31"},
		{name: "Slash30", input: "10.0.0.1 255.255.255.252", expected: "10.0.0.1/30"},
		{name: "Slash24", input: "10.0.0.1/24", expected: "10.0.0.1/24"},
		{name: "Slash22", input: "10.0.0.1/22", expected: "10.0.0.1/22"},
		{name: "Slash30Net", input: "10.0.0.0/30", expected: "10.0.0.0/30"},
		{name: "Slash30First", input: "10.0.0.1/30", expected: "10.0.0.1/30"},
		{name: "Slash30Last", input: "10.0.0.2/30", expected: "10.0.0.2/30"},
		{name: "Slash30Bcast", input: "10.0.0.3/30", expected: "10.0.0.3/30"},
		{name: "Slash1", input: "10.0.0.1/1", expected: "10.0.0.1/1"},
		{name: "Slash0", input: "10.0.0.1/0", expected: "10.0.0.1/0"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ipv4, err := ip.ParseIPv4(tc.input)

			// Check for unexpected error
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ipv4.String() != tc.expected {
				t.Errorf("expected string %q, got %q", tc.expected, ipv4.String())
			}
		})
	}
}

func TestIPv4UsableHosts(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected uint32
	}{
		{name: "Slash32", input: "10.0.0.1/32", expected: 0},
		{name: "Slash31", input: "10.0.0.1/31", expected: 2},
		{name: "Slash30", input: "10.0.0.1 255.255.255.252", expected: 2},
		{name: "Slash24", input: "10.0.0.1/24", expected: 254},
		{name: "Slash22", input: "10.0.0.1/22", expected: 1022},
		{name: "Slash30Net", input: "10.0.0.0/30", expected: 2},
		{name: "Slash30First", input: "10.0.0.1/30", expected: 2},
		{name: "Slash30Last", input: "10.0.0.2/30", expected: 2},
		{name: "Slash30Bcast", input: "10.0.0.3/30", expected: 2},
		{name: "Slash1", input: "10.0.0.1/1", expected: 2147483646},
		{name: "Slash0", input: "10.0.0.1/0", expected: 4294967294},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ipv4, err := ip.ParseIPv4(tc.input)

			// Check for unexpected error
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ipv4.UsableHosts() != tc.expected {
				t.Errorf("expected string %d, got %d", tc.expected, ipv4.UsableHosts())
			}
		})
	}
}

func TestIPv4NetworkSize(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected uint32
	}{
		{name: "Slash32", input: "10.0.0.1/32", expected: 1},
		{name: "Slash31", input: "10.0.0.1/31", expected: 2},
		{name: "Slash30", input: "10.0.0.1 255.255.255.252", expected: 4},
		{name: "Slash24", input: "10.0.0.1/24", expected: 256},
		{name: "Slash22", input: "10.0.0.1/22", expected: 1024},
		{name: "Slash30Net", input: "10.0.0.0/30", expected: 4},
		{name: "Slash30First", input: "10.0.0.1/30", expected: 4},
		{name: "Slash30Last", input: "10.0.0.2/30", expected: 4},
		{name: "Slash30Bcast", input: "10.0.0.3/30", expected: 4},
		{name: "Slash1", input: "10.0.0.1/1", expected: 2147483648},
		{name: "Slash0", input: "0.0.0.0/0", expected: 4294967295},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ipv4, err := ip.ParseIPv4(tc.input)

			// Check for unexpected error
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ipv4.NetworkSize() != tc.expected {
				t.Errorf("expected string %d, got %d", tc.expected, ipv4.NetworkSize())
			}
		})
	}
}
