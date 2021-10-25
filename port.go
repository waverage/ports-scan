package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"port/db"
	"port/ipslist"
	"port/limiter"
	"time"
)

var mongo = db.Db{}

func main() {
	gen, err := ipslist.Generator(0, 0, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = mongo.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Disconnect()

	fmt.Println("Database connected!")

	count := 0
	limit := limiter.NewConcurrencyLimiter(800)

	for ;; {
		ip, err := gen.Next()
		if err != nil {
			break
		}

		limit.Execute(func() {
			processPort(ip, "80")
		})

		if count % 1000000 == 0 {
			fmt.Println("--------------Processed ", count, ", IP: ", ip)
		}
		count++
	}

	limit.Wait()

	fmt.Printf("Total count: %d\n", count)
}

func processPort(host, port string) {
	err := checkPort(host, port)

	if err != nil {
		if err.Error() != fmt.Sprintf("dial tcp %v:%v: i/o timeout", host, port) &&
		   err.Error() != fmt.Sprintf("dial tcp %v:%v: connect: connection refused", host, port) {
			log.Println(host, ":", port, "Received error: ", err)
		}
	} else {
		log.Println(host + ":" + port + " opened")
		mongo.InsertRow(host, port)
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