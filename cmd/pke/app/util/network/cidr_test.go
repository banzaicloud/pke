// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		cidr     string
		ip       net.IP
		contains bool
		err      bool
	}{
		{"IP in CIDR", "10.240.0.0/24", net.ParseIP("10.240.0.1"), true, false},
		{"IP not in CIDR", "10.240.0.0/24", net.ParseIP("1.1.1.1"), false, false},
		{"invalid CIDR", "10.240.0.0/100", net.ParseIP("1.1.1.1"), false, true},
	}
	for _, tc := range testCases {
		contains, err := Contains(tc.cidr, tc.ip)
		if tc.err {
			require.Error(t, err, tc.name)
			continue
		} else {
			require.NoError(t, err, tc.name)
		}
		require.Equal(t, tc.contains, contains, tc.name)
	}
}

func TestContainsFirst(t *testing.T) {
	testCases := []struct {
		name     string
		cidr     string
		ips      []net.IP
		expected net.IP
		err      bool
	}{
		{"IP in CIDR", "10.240.0.0/24", []net.IP{net.ParseIP("10.240.0.1")}, net.ParseIP("10.240.0.1"), false},
		{"IP not in CIDR", "10.240.0.0/24", []net.IP{net.ParseIP("1.1.1.1")}, net.ParseIP("1.1.1.1"), true},
		{"invalid CIDR", "10.240.0.0/100", []net.IP{net.ParseIP("10.240.0.1")}, nil, true},
		{"multiple IPs in CIDR, first match", "10.240.0.0/24", []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("10.240.10.1"), net.ParseIP("10.240.0.1")}, net.ParseIP("10.240.0.1"), false},
	}
	for _, tc := range testCases {
		ip, err := ContainsFirst(tc.cidr, tc.ips)
		if tc.err {
			require.Error(t, err, tc.cidr, tc.ips)
			continue
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tc.expected, ip)
	}
}
