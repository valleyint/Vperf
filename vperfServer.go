package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const usualPort = 2222
const frameSize = 100 * 1024 * 1024

type server struct {
	conn *net.TCPConn
}

type frame struct {
	data []byte
}

type client struct {
	conn *net.TCPConn
}

func listen() (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", usualPort))
	if err != nil {
		return nil, err
	}

	listner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return listner, nil
}

func accept(listner *net.TCPListener) (*server, error) {
	conn, err := listner.AcceptTCP()
	if err != nil {
		_ = listner.Close()
		return nil, err
	}

	serv := server{conn: conn}
	return &serv, nil
}

func (s *server) flood() {

}

// add client here

func newFrame(size uint64) *frame {
	sendArr := make([]byte, size)
	_, _ = rand.Read(sendArr[9:])
	frm := frame{data: sendArr}
	frm.setSize(size - 9)
	return &frm
}

func (f *frame) setFinal(final bool) {
	if final {
		f.data[8] = 1
	} else {
		f.data[8] = 0
	}
}

func (f *frame) getFinal() bool {
	return f.data[8] == 1
}

func (f *frame) getSize() uint64 {
	return binary.BigEndian.Uint64(f.data[0:8])
}

func (f *frame) setSize(size uint64) {
	binary.BigEndian.PutUint64(f.data[0:8], size)
}

func doFlood(conn *net.TCPConn) error {
	frm := newFrame(frameSize)
	starTim := time.Now()
	for {
		last := time.Now().Sub(starTim) > 5*time.Second
		err := frm.send(conn, last)
		if err != nil {
			return err
		}

		if last {
			return nil
		}
	}
}

func (f *frame) send(conn *net.TCPConn, final bool) error {
	amountWriten := 0
	f.setFinal(final)
	for {
		numWriten, err := conn.Write(f.data[amountWriten:])
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

func (f *frame) receive(conn *net.TCPConn) error {
	var amountRead uint64
	for {
		numRead, err := conn.Read(f.data[amountRead:9])
		if err != nil {
			return err
		}

		amountRead += uint64(numRead)
		if amountRead == 9 {
			break
		}

		if numRead == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}

	amountRead = 0
	size := f.getSize()
	for {
		numRead, err := conn.Read(f.data[9 : size+9])
		if err != nil {
			return err
		}

		amountRead += uint64(numRead)
		if amountRead == size {
			return nil
		}

		if numRead == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}
}
