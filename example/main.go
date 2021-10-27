package main

import (
	"log"
	"net"
	"port/ipcore"
)

func main() {
	ip := net.IPv4(0, 0, 0, 0)

	for i := 0; i < 10; i++ {
		ip = ip.IpIncrement(ip, 256)
		log.Println("IP: ", ip.String())
	}

	// Generators test
	startIp := net.IPv4(0, 0, 0, 0)
	endIp := net.IPv4(0, 0, 0, 255)
	gen, err := ip.Generator(startIp, endIp)
	if err != nil {
		log.Fatal(err)
	}

	for ;; {
		ip, err := gen.Next()
		if err != nil {
			break
		}
		log.Println("IP: ", ip)
	}
}
