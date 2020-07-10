package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const Port = 2222
const Size = 100 * 1024 * 1024

type server struct {
	conn *net.TCPConn
}

type frame struct {
	data []byte
}

type client struct {
	conn *net.TCPConn
}
func listen () (*net.TCPListener , error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", Port))
	if err != nil {
		return nil, err
	}

	listner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return listner , nil
}

func accept (listner *net.TCPListener) (*server , error){
	conn , err := listner.AcceptTCP()
	if err != nil {
		_ = listner.Close()
		return nil , err
	}

	serv := server{conn: conn}
	return &serv , nil
}

func (s *server) flood () {}

// add client here

func newFrame () *frame {
	sendArr := make([]byte , Size)
	_ , _ = rand.Read(sendArr[9 : ])
	binary.BigEndian.PutUint64(sendArr[0 : 8] , uint64(len(sendArr)))

	frm := frame{data: sendArr}
	return &frm
}

func (f *frame) finalise (final bool) {
	if final {
		f.data[8] = 1
	} else {
		f.data[8] = 0
	}
}

func doFlood (conn *net.TCPConn) (int , error) {
	frm := newFrame()
	starTim := time.Now()
	loop := 0
	for loop = 0 ; ; loop ++ {
		totalWriten := 0
		for {
			numWriten, err := conn.Write(frm.data)
			if err != nil {
				return loop , err
			}
			totalWriten = totalWriten + numWriten

			if totalWriten == Size {
				break
			}
		}

		if time.Now().Sub(starTim) >= 5 {
			break
		}
	}

	return loop , nil
}

func (f *frame) send (conn *net.TCPConn , final bool) error {
	amountWriten := 0
	f.finalise(final)
	for {
		numWriten, err := conn.Write(f.data[amountWriten : ])
		if err != nil {
			return err
		}

		amountWriten += numWriten
		if amountWriten == len(f.data) {
			return nil
		}

		if numWriten == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}
}
