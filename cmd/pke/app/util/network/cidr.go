package network

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

func Contains(cidr string, ip net.IP) (bool, error) {
	_, c, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	return c.Contains(ip), nil
}

func ContainsFirst(cidr string, ips []net.IP) (net.IP, error) {
	for _, ip := range ips {
		if ok, err := Contains(cidr, ip); err == nil && ok {
			return ip, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cidr %q does not contain ip %q", cidr, ips))
}
