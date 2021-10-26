package ipslist

import (
	"errors"
	"net"
)

type Ip struct {
	a, b, c, d byte
}

type Ips struct {
	cursor net.IP
	endIp net.IP
}

func (f *Ips) Next() (string, error) {
	result := f.cursor.String()
	if f.cursor.String() == f.endIp.String() {
		return "", errors.New("all ips returned")
	} else {
		f.cursor = IpIncrement(f.cursor, 1)
	}

	return result, nil
}

func Generator(startIp net.IP, endIp net.IP) (Ips, error) {
	return Ips{startIp, endIp}, nil
}

func IpIncrement(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}