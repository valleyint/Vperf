package main

import (
	"net"
	"fmt"
)

const Port = 2222

type server struct {
	conn *net.TCPConn
}

func listuin () (*server , error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", Port))
	if err != nil {
		return nil, err
	}

	listner , err := net.ListenTCP("tcp" , addr)
	if err != nil {
		return nil , err
	}

	conn , err := listner.AcceptTCP()
	if err != nil {
		listner.Close()
		return nil , err
	}

	serv := server{conn: conn}
	return &serv , nil
}


