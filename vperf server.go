package main

import (
	"fmt"
	"net"
)

const Port = 2222

type server struct {
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

func accsept (listner *net.TCPListener) (*server , error){
	conn , err := listner.AcceptTCP()
	if err != nil {
		_ = listner.Close()
		return nil , err
	}

	serv := server{conn: conn}
	return &serv , nil
}
