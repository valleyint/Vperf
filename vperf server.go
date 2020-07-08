package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

const Port = 2222
const Size = 100 * 1024 * 1024

type server struct {
	conn *net.TCPConn
}

type frame struct {
	buff []byte
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

func accsept (listner *net.TCPListener) (*server , error){
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

func NewFrame () *frame {
	sendArr := make([]byte , Size - 8)
	_ , _ = rand.Read(sendArr)
	binary.BigEndian.PutUint64(sendArr[0 : 8] , uint64(len(sendArr)))
	frm := frame{buff: sendArr}

	return &frm
}

func doFlood (conn *net.TCPConn) {
	frm := NewFrame()

	
}