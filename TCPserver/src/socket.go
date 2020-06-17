package main

import (
	"log"
	"net"
	"os"
	"time"
)

type Socket struct {
	TCPAddr *net.TCPAddr
	Handler *net.TCPListener
}

func (s *Socket) SetTcpAddr() {
	var err error
	(*s).TCPAddr, err = net.ResolveTCPAddr("tcp", os.Getenv("TCPConn"))
	if err != nil {
		log.Printf("TCPserver, TCPAddrSet: %v", err)
	}
}

func (s *Socket) setTcpHandler() {
	var err error
	s.Handler, err = net.ListenTCP("tcp", s.TCPAddr)
	if err != nil {
		log.Printf("TCPserver, setTcpHandler: %v", err)
	}

	log.Println("Server started on the :", s.TCPAddr.String())
}

// setDeadline - (t int) t = Millisecond
func (s *Socket) setDeadline(t int) {
	var DeadlineTime = time.Now().Add(time.Duration(t))
	err := s.Handler.SetDeadline(DeadlineTime)
	if err != nil {
		log.Printf("TCPserver, setDeadline: %v", err)
	}
}

func (s *Socket) acceptConn(conns *conns) {
	for {
		conn, err := s.Handler.AcceptTCP()
		if err != nil {
			log.Printf("TCPserver, acceptConn: %v", err)
		} else {
			err := conns.addConn(conn)
			if err != nil {
				log.Printf("TCPserver, acceptConn: %v", err)
			}
		}
	}
}
