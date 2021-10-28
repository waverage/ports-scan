package main

import (
	"errors"
	"net"
	"time"
)

func CheckPort(host string, port string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return err
	}

	if conn != nil {
		conn.Close()
		return nil
	}

	return errors.New("connection failed")
}

