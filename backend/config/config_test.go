package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGenerateIPRanges(t *testing.T) {
	testCases := []struct {
		serverIP string
		mask     string
		expected int
		err      string
	}{
		{"192.168.1.10", "/24", 2, ""},
		{"192.168.1.1", "/24", 1, ""},
		{"192.168.1.254", "/24", 1, ""},
		{"192.168.1.1", "/30", 1, ""},
		// {"192.168.2.1", "/24", 0, ""}, // TODO
		{"192.168.1.0", "/24", 0, "server IP is the network address"},
		{"192.168.1.255", "/24", 0, "server IP is the broadcast address"},
		{"192.168.10.50", "/16", 257, ""},
		{"192.168.1.1", "/32", 0, "server IP is the network address"},
		{"192.168.1.1", "/31", 0, "server IP is the broadcast address"},
		{"192.168.1.1", "/xyz", 0, "invalid CIDR address: 192.168.1.1/xyz"},
		{"invalid_ip", "/24", 0, "invalid CIDR address: invalid_ip/24"},
		{"192.168.1.3", "/29", 2, ""},
		// {"10.0.0.1", "/8", 65729, ""}, // not working??
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			result, err := GenerateIPRanges(tc.serverIP, tc.mask)
			if tc.err != "" {
				if err == nil || err.Error() != tc.err {
					t.Errorf("expected error: %s, got: %v", tc.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(len(result), tc.expected) {
					t.Errorf("expected: %+v, got: %+v, %+v", tc.expected, len(result), 0)
				}
			}
		})
	}
}
