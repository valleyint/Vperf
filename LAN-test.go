package main

import (
	"net"
	"time"
)

const aces = "aces.lan.melkote.com"

func DPerfClient () (int , error) {
	addr , err := net.ResolveTCPAddr("tcp" , aces)
	if err != nil {
		return 0 , err
	}

	conn , err := net.DialTCP("tcp" , nil , addr)
	if err != nil {
		return 0 , err
	}

	buff := make([]byte , 1000000)

	startTim := time.Now()
	for loop := 0 ; loop < 1000 ; loop ++ {
		_ , err = conn.Write(buff)
		if err != nil {
			return 0 , err
		}
	}
	endTim := time.Now()

	tim := endTim.Sub(startTim)
	return int(tim.Milliseconds()) , nil
}
 