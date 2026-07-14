package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mukeshmahato17/goredis/resp"
)

type Server struct {
	listenAddr      string
	connectedClient int
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() {
	fmt.Println("server started")
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	s.acceptLoop(ln)
}

func (s *Server) acceptLoop(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		s.connectedClient++
		fmt.Println("client connected:", conn.RemoteAddr(), "connected client:", s.connectedClient)
		s.readConn(conn)
	}
}

func (s *Server) readConn(conn net.Conn) {
	buf := make([]byte, 512)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.connectedClient--
				fmt.Println("client disconnected:", conn.RemoteAddr(), "connected client:", s.connectedClient)
				break
			}
			fmt.Println("error", err)
		}
		val, err := resp.Decode(buf)
		fmt.Println(val)

		conn.Write([]byte("+OK\r\n"))
	}
}
