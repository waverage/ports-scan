package main

import (
	"fmt"
	"log"
	"net"
	"port/db"
	"port/ipcore"
	"sync"
	"sync/atomic"
)

var counter uint64
var wg sync.WaitGroup
var mg = db.Db{}

const useDb = false
const totalIps = uint(256 * 256 * 256 * 256)

var reservedNetworks []ipcore.NetworkRangeInt

func main() {
	if useDb {
		err := mg.Connect()
		if err != nil {
			log.Fatal(err)
		}
		defer mg.Disconnect()
		fmt.Println("Database connected!")
	}

	reservedNetworks = ipcore.GetReservedNetworks()

	totalWorkers := uint(4096)
	ipsPerWorker := uint(totalIps / totalWorkers)

	log.Println("total ips: ", totalIps)
	log.Println("Ips per worker: ", ipsPerWorker)
	ip := net.IPv4(0, 0, 0, 0)

	processedIps := uint64(0)

	wg.Add(int(totalWorkers))

	for i := uint(1); i <= totalWorkers; i++ {
		var inc uint
		if i == totalWorkers {
			// Fix last network range (divide error fix)
			inc = uint(uint64(totalIps) - processedIps - 1)
		} else {
			inc = ipsPerWorker
		}

		endIp := ipcore.Increment(ip, inc)
		log.Println("Worker", i, "network", ip.String(), "-", endIp.String())

		go doWorker(ip, inc, i)

		ip = ipcore.Increment(ip, ipsPerWorker)
		processedIps += uint64(ipsPerWorker)
	}

	wg.Wait()

	log.Println("All workers has been initialized")
}

func doWorker(ip net.IP, count uint, index uint) {
	defer wg.Done()

	//log.Println("#", index, " initial ip:", ip.String())

	for i := uint(0); i < count; i++ {
		if ipcore.IpIsReserved(ip, reservedNetworks) {
			// Skip reserved network
			nToSkip := ipcore.GetNToEndOfReservedNetwork(ip, reservedNetworks)
			if nToSkip > 0 {
				newIp := ipcore.Increment(ip, uint(nToSkip))
				log.Println("#", index, "Skip", nToSkip, " reserved ips", ", reserved ip:", ip.String(), ", end reserved network: ", newIp)
				ip = newIp

				if ip.String() == "0.0.0.0" {
					log.Println("#", index, "group finished, because all ips are in reserved network")
					return
				}
			}
		}

		processPort(ip, "80")
		atomic.AddUint64(&counter, 1)

		ip = ipcore.Increment(ip, 1)

		if counter % 100000 == 0 {
			log.Println("#", index, "Processed ", counter / 1000000, "mln ips, ", ip.String())
		}
	}

	log.Println("#", index, "Group finished, end:", ip.String())
}

func processPort(host net.IP, port string) {
	// Skip reserved ips
	if ipcore.IpIsReserved(host, reservedNetworks) {
		return
	}

	stringIp := host.String()
	err := CheckPort(stringIp, port)

	if err != nil {
		if err.Error() != fmt.Sprintf("dial tcp %v:%v: i/o timeout", host, port) &&
			err.Error() != fmt.Sprintf("dial tcp %v:%v: connect: network is unreachable", host, port) &&
			err.Error() != fmt.Sprintf("dial tcp %v:%v: connect: no route to host", host, port) &&
		    err.Error() != fmt.Sprintf("dial tcp %v:%v: connect: connection refused", host, port) {
			log.Println(host, ":", port, "Received error: ", err)
		}
	} else {
		if useDb {
			mg.InsertRow(stringIp, port)
		}
	}
}
