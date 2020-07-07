package main

import (
	"net"
)

type server struct {
	addr string
	conn *net.TCPConn
}

func listuin () (*server , error) {
	addr, err := net.ResolveTCPAddr("tcp", "126.0.0.1:2222")
	if err != nil {
		return nil, err
	}

	listner , err := net.ListenTCP("tcp" , addr)
	if err != nil {
		return nil , err
	}

	conn , err := listner.AcceptTCP()
	if err != nil {
		return nil , err
	}

	serv := server{conn: conn , addr: conn.RemoteAddr().String()}
	return &serv , nil
}


