package config

import (
	"net"
	"testing"
)

func TestGenerateIPRanges(t *testing.T) {
	tests := []struct {
		name      string
		serverIP  string
		mask      string
		expected  []Ip_Range
		expectErr bool
	}{
		// Valid Inputs
		{
			name:     "Valid small range",
			serverIP: "192.168.1.0",
			mask:     "/30",
			expected: []Ip_Range{
				{Start: net.ParseIP("192.168.1.1").To4(), End: net.ParseIP("192.168.1.2").To4()},
			},
			expectErr: false,
		},
		{
			name:     "Valid /29 subnet",
			serverIP: "10.0.0.0",
			mask:     "/29",
			expected: []Ip_Range{
				{Start: net.ParseIP("10.0.0.1").To4(), End: net.ParseIP("10.0.0.6").To4()},
			},
			expectErr: false,
		},

		// Edge Cases
		{
			name:      "Invalid CIDR input",
			serverIP:  "192.168.1.1",
			mask:      "/33",
			expectErr: true,
		},
		{
			name:      "Invalid IP format",
			serverIP:  "invalid_ip",
			mask:      "/24",
			expectErr: true,
		},
		{
			name:     "Exact range size match",
			serverIP: "10.0.0.0",
			mask:     "/28",
			expected: []Ip_Range{
				{Start: net.ParseIP("10.0.0.1").To4(), End: net.ParseIP("10.0.0.14").To4()},
			},
			expectErr: false,
		},
		{
			name:      "Single IP range 32",
			serverIP:  "10.0.0.1",
			mask:      "/32",
			expected:  []Ip_Range{},
			expectErr: false,
		},

		// Special IP Addresses
		{
			name:     "Broadcast and network exclusion",
			serverIP: "192.168.0.0",
			mask:     "/29",
			expected: []Ip_Range{
				{Start: net.ParseIP("192.168.0.1").To4(), End: net.ParseIP("192.168.0.6").To4()},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranges, err := GenerateIPRanges(tt.serverIP, tt.mask)
			if (err != nil) != tt.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.expectErr {
				// Validate ranges
				if len(ranges) != len(tt.expected) {
					t.Fatalf("expected %d ranges, got %d", len(tt.expected), len(ranges))
				}

				for i, r := range ranges {
					expectedStart := tt.expected[i].Start
					expectedEnd := tt.expected[i].End

					if !net.IP.Equal(r.Start, expectedStart) || !net.IP.Equal(r.End, expectedEnd) {
						t.Errorf("range %d mismatch: got start %v, end %v; expected start %v, end %v",
							i, r.Start, r.End, expectedStart, expectedEnd)
					}
				}
			}
		})
	}
}
