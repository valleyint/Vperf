package main

import (
	"net"
	"os"
	"time"
)


func DPerfClient() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", aces)
	if err != nil {
		return 0, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return 0, err
	}

	buff := make([]byte, 1024*1024)

	startTim := time.Now()
	for loop := 0; loop < 1024; loop++ {
		_, err = conn.Write(buff)
		if err != nil {
			return 0, err
		}

		if loop%10 == 0 && loop != 0 {
			readable := []byte{}
			_, err = conn.Read(readable)
			if err != nil {
				return 0, err
			}

		}
	}
	endTim := time.Now()

	err = conn.Close()
	if err != nil {
		return 0, err
	}
	tim := endTim.Sub(startTim)
	return int(tim.Milliseconds()), nil
}

// /var/server should be less than 100 bytes long
func findServer() error {
	file, err := os.Open("/var/server")
	if err != nil {
		return err
	}

	readable := make([]byte, 100)
	_, err = file.Read(readable)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
