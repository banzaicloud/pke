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

package kubeadm

import (
	"fmt"
	"math"
	"net"
)

func SplitHostPort(hostport, defaultPort string) (host, port string, err error) {
	host, port, err = net.SplitHostPort(hostport)
	if aerr, ok := err.(*net.AddrError); ok {
		if aerr.Err == "missing port in address" {
			hostport = net.JoinHostPort(hostport, defaultPort)
			host, port, err = net.SplitHostPort(hostport)
		}
	}
	return
}

func KubeReservedMemory(totalMemoryBytes uint64) string {
	// convert bytes to MiB
	mem := totalMemoryBytes >> 20

	var multiplier float64

	switch {
	case mem > 200000:
		multiplier = 0.06
	case mem > 100000:
		multiplier = 0.076
	case mem > 50000:
		multiplier = 0.092
	case mem > 25000:
		multiplier = 0.123
	case mem > 15000:
		multiplier = 0.166
	case mem > 7500:
		multiplier = 0.225
	case mem > 3750:
		multiplier = 0.243
	case mem > 1700:
		multiplier = 0.256
	case mem > 700:
		multiplier = 0.3
	case mem <= 700 && mem > 100:
		multiplier = 0.4
	default:
		multiplier = 0
	}

	return fmt.Sprintf("%dMi", uint64(math.Round(float64(mem)*multiplier)))
}
