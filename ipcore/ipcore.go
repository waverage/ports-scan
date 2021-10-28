package ipcore

import (
	"encoding/binary"
	"errors"
	"net"
)

func Increment(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}

func Ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func Int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

type NetworkRangeIp struct {
	a, b net.IP
}

type NetworkRangeInt struct {
	a, b uint32
}

// ReservedNetworks https://en.wikipedia.org/wiki/Reserved_IP_addresses
var ReservedNetworks = []NetworkRangeIp{
	{net.IPv4(0, 0, 0, 0), net.IPv4(0, 255, 255, 255)},
	{net.IPv4(10, 0, 0, 0), net.IPv4(10, 255, 255, 255)},
	{net.IPv4(100, 64, 0, 0), net.IPv4(100, 127, 255, 255)},
	{net.IPv4(127, 0, 0, 0), net.IPv4(127, 255, 255, 255)},
	{net.IPv4(169, 254, 0, 0), net.IPv4(169, 254, 255, 255)},
	{net.IPv4(172, 16, 0, 0), net.IPv4(172, 31, 255, 255)},
	{net.IPv4(192, 0, 0, 0), net.IPv4(192, 0, 0, 255)},
	{net.IPv4(192, 0, 2, 0), net.IPv4(192, 0, 2, 255)},
	{net.IPv4(192, 88, 99, 0), net.IPv4(192, 88, 99, 255)},
	{net.IPv4(192, 168, 0, 0), net.IPv4(192, 168, 255, 255)},
	{net.IPv4(198, 18, 0, 0), net.IPv4(198, 19, 255, 255)},
	{net.IPv4(198, 51, 100, 0), net.IPv4(198, 51, 100, 255)},
	{net.IPv4(203, 0, 113, 0), net.IPv4(203, 0, 113, 255)},
	{net.IPv4(224, 0, 0, 0), net.IPv4(239, 255, 255, 255)},
	{net.IPv4(240, 0, 0, 0), net.IPv4(255, 255, 255, 255)},
}

func GetReservedNetworks() []NetworkRangeInt {
	result := make([]NetworkRangeInt, len(ReservedNetworks))

	for index, network := range ReservedNetworks {
		result[index] = NetworkRangeInt{Ip2int(network.a), Ip2int(network.b)}
	}

	return result
}

func IpIsReserved(ip net.IP, reserved []NetworkRangeInt) bool {
	intIp := Ip2int(ip)

	allReservedAfter := Ip2int(net.IPv4(224, 0, 0, 0))
	if intIp >= allReservedAfter {
		return true
	}

	for _, network := range reserved {
		if intIp >= network.a && intIp <= network.b {
			return true
		}
	}

	return false
}

func GetNToEndOfReservedNetwork(ip net.IP, reserved []NetworkRangeInt) (error, uint32) {
	intIp := Ip2int(ip)

	allReservedAfter := Ip2int(net.IPv4(224, 0, 0, 0))
	if intIp >= allReservedAfter {
		return errors.New("do not have free ips after"), 0
	}

	for _, network := range reserved {
		if intIp >= network.a && intIp <= network.b {
			return nil, (network.b - intIp) + 1
		}
	}

	return nil, 0
}

func SkipReserved(ip net.IP, limit uint, reserved []NetworkRangeInt) (error, net.IP, uint) {
	skipped := uint(0)

	for true {
		if IpIsReserved(ip, reserved) {
			err, skip := GetNToEndOfReservedNetwork(ip, reserved)
			if err != nil {
				return err, nil, 0
			}

			skipped += uint(skip)
			if skipped >= limit || ip.String() == "0.0.0.0" {
				return errors.New("all ips in network are reserved"), nil, 0
			}

			ip = Increment(ip, uint(skip))
		} else {
			break
		}
	}

	return nil, ip, limit - skipped
}