package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const usualPort = "2222"
const maxBufferSize = 100 * 1024 * 1024

type options struct {
	server bool
	client string
}

type server struct {
	conn *net.TCPConn
}

type frame struct {
	data []byte
	firstArrival time.Time
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

func (s *server) flood () (string , error) {
	frm := newFrame(maxBufferSize)
	//the only command supported is flood me , so we ignore cmd frame
	err := frm.receive(s.conn)
	if err != nil {
		return "" , nil
	}
	fmt.Println("recived order")
	err = doFlood(s.conn , frm)
	if err != nil {
		return "" , err
	}

	err = frm.receive(s.conn)
	if err != nil {
		return "" , err
	}

	stat := string(frm.data[9 : frm.getSize()])

	return stat , nil
}

// add client here

func newFrame(size uint64) *frame {
	sendArr := make([]byte, size)
	_, _ = rand.Read(sendArr[9:])
	frm := frame{data: sendArr}
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

func doFlood(conn *net.TCPConn , frm *frame) error {
	starTim := time.Now()
	frm.setSize(maxBufferSize)
	for {
		fmt.Println("sentframe")
		last := time.Now().Sub(starTim) >= 5 * time.Second
		err := frm.send(conn, last)
		if err != nil {
			return err
		}

		if last {
			return nil
		}
	}
}

//latency milliseconds, throughput Megabits/second , error
func handleFlood(conn *net.TCPConn , frm *frame) (int , float64 , error) {
	cmd := []byte("flood me")
	copy(frm.data[9 : ] , cmd)
	frm.setSize(uint64(len(cmd)+9))

	Start := time.Now()
	err := frm.send(conn , true)
	if err != nil {
		return 0 , 0 ,err
	}

	bytesRead := uint64(0)
	lat := 0
	for loop := 0 ; ; loop ++ {
		err = frm.receive(conn)
		if err != nil {
			return 0 , 0 , err
		}

		if loop == 0 {
			lat = int(frm.firstArrival.Sub(Start).Milliseconds())
		}

		if frm.getFinal() == true {
			fmt.Println("Now received final", frm.getFinal())
			break
		}

		bytesRead += frm.getSize()
		fmt.Println("Now received total flood", bytesRead)
	}
	throughDur := time.Now().Sub(Start)
	throughPut := float64(bytesRead) / throughDur.Seconds()
	throughPut = (throughPut *8) / (1024 * 1024)

	return lat , throughPut , nil
}

func (f *frame) send(conn *net.TCPConn, final bool) error {
	amountWriten := 0
	f.setFinal(final)
	sz := f.getSize()
	fmt.Println("Will send total of", sz, "bytes", "final", final)
	for {
		numWriten, err := conn.Write(f.data[amountWriten : sz])
		if err != nil {
			return err
		}

		amountWriten += numWriten
		if amountWriten == int(sz) {
			return nil
		}

		fmt.Println("Now written total of", amountWriten)

		if numWriten == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (f *frame) receive (conn *net.TCPConn) error {
	var amountRead uint64

	for amountRead < 1 {
		numRead, err := conn.Read(f.data[0:1])
		if err != nil {
			return err
		}
		amountRead += uint64(numRead)

		if numRead == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	f.firstArrival = time.Now()
	fmt.Println("Got first byte")

	for {
		numRead, err := conn.Read(f.data[amountRead:9])
		if err != nil {
			return err
		}

		amountRead += uint64(numRead)
		fmt.Println("Total read now is", amountRead)
		if amountRead == 9 {
			break
		}

		if numRead == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}

	size := f.getSize()
	fmt.Println("Now reading remaining", size)
	for {
		numRead, err := conn.Read(f.data[amountRead : size])
		if err != nil {
			return err
		}

		amountRead += uint64(numRead)
		fmt.Println("Now read total", amountRead)
		if amountRead == size {
			fmt.Println("Done reading", amountRead)
			return nil
		}

		if numRead == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func connect (host string) (*client , error) {
	raddr , err := net.ResolveTCPAddr("tcp" , net.JoinHostPort(host , usualPort))
	if err != nil {
		return nil, err
	}

	conn , err := net.DialTCP("tcp" , nil , raddr)
	if err != nil {
		return nil, err
	}

	clint := client{conn: conn}
	return &clint , nil
}

func (c *client) askflood () (int , float64 , error) {
	frm := newFrame(maxBufferSize)
	lat , speed , err := handleFlood(c.conn , frm)
	if err != nil {
		return 0 , 0 , err
	}

	stat := fmt.Sprintf("Done askFlood latency was %v ms speed %v mbps\n", lat, speed)
	copy(frm.data[9 : ] , stat)
	frm.setSize(uint64(len(stat)+9))

	err = frm.send(c.conn , true)
	if err != nil {
		return 0 , 0 ,err
	}
	fmt.Println("latancy" , lat , "speed" , speed)
	return lat , speed , nil
}

func Vperf () {
	var opts options
	flag.BoolVar(&opts.server , "server" , false , "specify a mode -s (server) or -c (client)")
	flag.StringVar(&opts.client , "client" , "" , "specify a mode -s (server) or -c (client)")
	flag.Parse()

	if opts.server && len(opts.client) != 0 {
		fmt.Printf("the program cannot be both a sever and client /n if you wish to loopback please use 2 windows for the function")
		return
	}

	if !opts.server && len(opts.client) == 0 {
		fmt.Println("specify a mode -s (server) or -c (client)")
		return
	}

	if opts.server {
		listner , err := listen()
		if err != nil {
			panic(err)
		}
		fmt.Println("listning")

		serv , err := accept(listner)
		if err != nil {
			panic(err)
		}
		fmt.Println("serving")

		fmt.Println("floodstart")
		stat , err := serv.flood()
		if err != nil {
			panic(err)
		}

		fmt.Println(stat)
	} else {
		clint , err := connect(opts.client)
		if err != nil {
			panic(err)
		}
		fmt.Println("connected")

		fmt.Println("askedFlood")
		lat , sped , err := clint.askflood()
		if err != nil {
			panic(err)
		}
		fmt.Println(lat , sped)
	}
}

func main () {
	Vperf()
}
