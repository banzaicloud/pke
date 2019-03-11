package network

import (
	"net"

	"github.com/pkg/errors"
)

func IPv4Addresses() ([]net.IP, error) {
	ips := make([]net.IP, 0)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip)
		}
	}

	if len(ips) == 0 {
		return nil, errors.New("no IPv4 address found")
	}

	return ips, nil
}
