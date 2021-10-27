package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"port/db"
	"port/ipslist"
	"sync"
	"sync/atomic"
	"time"
)


var counter uint64
var wg sync.WaitGroup
var mg = db.Db{}

func main() {
	err := mg.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer mg.Disconnect()

	fmt.Println("Database connected!")

	totalIps := 255 * 255 * 255 * 255
	totalWorkers := 16
	ipsPerWorker := uint(totalIps / totalWorkers)

	log.Println("total ips: ", totalIps)
	log.Println("Ips per worker: ", ipsPerWorker)
	ip := net.IPv4(0, 0, 0, 0)


	for i := 0; i <= totalWorkers; i++ {
		wg.Add(1)

		if i == totalWorkers {
			ip = net.IPv4(255, 255, 255, 255)
		}
		//if i % 1000 == 0 {
			log.Println("Start ip: ", ip.String())
		//}
		go initWorker(ip, ipsPerWorker)

		ip = ipslist.IpIncrement(ip, ipsPerWorker)
	}

	wg.Wait()

	log.Println("Workers has been initialized")
}

func initWorker(ip net.IP, count uint) {
	defer wg.Done()

	for i := uint(0); i < count; i++ {
		processPort(ip.String(), "80")
		atomic.AddUint64(&counter, 1)

		ip = ipslist.IpIncrement(ip, 1)

		if counter % 100000 == 0 {
			log.Println("Processed ", counter / 1000000, "mln ips, ", ip.String())
		}
	}

	log.Println("Group finished, end:", ip.String())
}

func processPort(host, port string) {
	err := checkPort(host, port)

	if err != nil {
		if err.Error() != fmt.Sprintf("dial tcp %v:%v: i/o timeout", host, port) &&
			err.Error() != fmt.Sprintf("dial tcp %v:%v: connect: network is unreachable", host, port) &&
			err.Error() != fmt.Sprintf("dial tcp %v:%v: connect: no route to host", host, port) &&
		    err.Error() != fmt.Sprintf("dial tcp %v:%v: connect: connection refused", host, port) {
			log.Println(host, ":", port, "Received error: ", err)
		}
	} else {
		mg.InsertRow(host, port)
	}
}

func checkPort(host string, port string) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Second / 5)
	if err != nil {
		return err
	}

	if conn != nil {
		conn.Close()
		return nil
	}

	return errors.New("connection failed")
}
