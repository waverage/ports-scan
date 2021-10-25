package ipslist

import (
	"errors"
	"net"
)

type Ips struct {
	ip net.IP
}

func (f *Ips) Next() (string, error) {
	result := f.ip.String()
	if f.ip.String() == "255.255.255.255" {
		return "", errors.New("all ips returned")
	} else {
		inc(f.ip)
	}

	return result, nil
}
//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Generator(a, b, c, d byte) (Ips, error) {
	return Ips{net.IPv4(a, b, c, d)}, nil
}