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
	for {
		cmd, err := readCommand(conn)
		if err != nil {
			conn.Close()
			s.connectedClient--
			fmt.Println("client disconnected:", conn.RemoteAddr(), "connected client:", s.connectedClient)
			if err == io.EOF {
				break
			}
		}
		respond(cmd, conn)
	}
}

func readCommand(conn net.Conn) (*resp.RedisCmd, error) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf[:])
	if err != nil {
		return nil, err
	}

	token, err := resp.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &resp.RedisCmd{
		Command: token[0],
		Args:    token[1:],
	}, nil
}

func respond(cmd *resp.RedisCmd, conn net.Conn) {
	fmt.Println(cmd)
	err := resp.EvalAndRespond(cmd, conn)
	if err != nil {
		respondError(err, conn)
	}
}

func respondError(err error, conn net.Conn) {
	conn.Write([]byte(fmt.Sprintf("%s\r\n", err)))
}
