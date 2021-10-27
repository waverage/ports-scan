package ipcore

import (
	"net"
	"testing"
)

var totalIps = uint(256 * 256 * 256 * 256)
var network24 = uint(256 * 256)
var network32 = uint(256)

type testpair struct {
	input net.IP
	inc uint
	expected net.IP
}

var tests = []testpair{
	{net.IPv4(0, 0, 0, 0), 10, net.IPv4(0, 0, 0, 10)},
	{net.IPv4(0, 0, 0, 0), 255, net.IPv4(0, 0, 0, 255)},
	{net.IPv4(0, 0, 0, 0), network32, net.IPv4(0, 0, 1, 0)},
	{net.IPv4(0, 0, 0, 0), network24, net.IPv4(0, 1, 0, 0)},
	{net.IPv4(0, 0, 0, 0), network24 * 50, net.IPv4(0, 50, 0, 0)},
	{net.IPv4(0, 0, 0, 0), totalIps - 1, net.IPv4(255, 255, 255, 255)},
}

func TestIpIncrement(t *testing.T) {
	for _, pair := range tests {
		v := Increment(pair.input, pair.inc)

		if v.String() != pair.expected.String() {
			t.Error(
				"For input:", pair.input.String(),
				"inc:", pair.inc,
				"expected:", pair.expected.String(),
				"got:", v.String())
		}
	}
}

type ip2IntPair struct {
	input net.IP
	expected uint32
}

var ip2IntTests = []ip2IntPair {
	{net.IPv4(0, 0, 0, 0), uint32(0)},
	{net.IPv4(0, 0, 1, 0), uint32(256)},
	{net.IPv4(0, 0, 255, 0), uint32(65280)},
	{net.IPv4(255, 255, 255, 255), uint32(totalIps - 1)},
}

func TestIp2int(t *testing.T) {
	for _, pair := range ip2IntTests {
		v := Ip2int(pair.input)

		if v != pair.expected {
			t.Error(
				"For input:", pair.input.String(),
				"expected:", pair.expected,
				"got:", v,
			)
		}
	}
}

type int2IpPair struct {
	input uint32
	expected net.IP
}

var int2IpTests = []int2IpPair {
	{uint32(0), net.IPv4(0, 0, 0, 0)},
	{uint32(256), net.IPv4(0, 0, 1, 0)},
	{uint32(65280), net.IPv4(0, 0, 255, 0)},
	{uint32(totalIps - 1), net.IPv4(255, 255, 255, 255)},
}

func TestInt2ip(t *testing.T) {
	for _, pair := range int2IpTests {
		v := Int2ip(pair.input)

		if !v.Equal(pair.expected) {
			t.Error(
				"For input:", pair.input,
				"expected:", pair.expected.String(),
				"got:", v,
			)
		}
	}
}

var reservedNetworks = GetReservedNetworks()

type IpIsReservedPair struct {
	input net.IP
	expected bool
}

var ipIsReservedTests = []IpIsReservedPair {
	{net.IPv4(0, 0, 0, 0), true},
	{net.IPv4(0, 0, 1, 0), true},
	{net.IPv4(0, 1, 0, 0), true},
	{net.IPv4(0, 255, 0, 0), true},
	{net.IPv4(0, 255, 255, 255), true},
	{net.IPv4(1, 0, 0, 0), false},

	{net.IPv4(9, 255, 255, 255), false},
	{net.IPv4(10, 0, 0, 0), true},
	{net.IPv4(10, 0, 0, 1), true},
	{net.IPv4(10, 255, 255, 255), true},
	{net.IPv4(11, 0, 0, 0), false},

	{net.IPv4(192, 168, 0, 0), true},
	{net.IPv4(192, 168, 255, 255), true},

	{net.IPv4(240, 0, 0, 0), true},
	{net.IPv4(240, 255, 0, 0), true},
	{net.IPv4(255, 255, 255, 255), true},
}

func TestIpIsReserved(t *testing.T) {
	for _, pair := range ipIsReservedTests {
		v := IpIsReserved(pair.input, reservedNetworks)

		if v != pair.expected {
			t.Error(
				"For input:", pair.input,
				"expected:", pair.expected,
				"got:", v,
			)
		}
	}
}


type endOfReservedPair struct {
	input net.IP
	expected uint32
}

var endOfReservedTests = []endOfReservedPair {
	{net.IPv4(0, 255, 255, 255), 1},
	{net.IPv4(0, 255, 255, 0), 256},
	{net.IPv4(192, 0, 0, 0), 256},
	{net.IPv4(192, 0, 0, 250), 6},

	// Return 0 for not reserved IP
	{net.IPv4(20, 10, 0, 0), 0},

	{net.IPv4(255, 255, 255, 255), 1},
}

func TestGetNToEndOfReservedNetwork(t *testing.T) {
	for _, pair := range endOfReservedTests {
		v := GetNToEndOfReservedNetwork(pair.input, reservedNetworks)

		if v != pair.expected {
			t.Error(
				"Get N to end of reserved network failed:",
				"For input:", pair.input,
				"expected:", pair.expected,
				"got:", v,
			)
		}
	}
}
