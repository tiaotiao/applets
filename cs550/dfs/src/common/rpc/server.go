package rpc

import (
	"common"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"common/log"
)

type Server struct {
	listener net.Listener
	handler  func(net.Conn) interface{}

	port    int
	stopped bool
	ch      chan error
}

func NewServer(port int, handler func(net.Conn) interface{}) *Server {
	s := &Server{}
	s.port = port
	s.handler = handler

	s.stopped = true
	s.ch = make(chan error)
	return s
}

func (s *Server) Stop() {
	if s.stopped {
		return
	}
	s.stopped = true

	s.listener.Close() // Close listener

	<-s.ch // Wait for signal
}

func (s *Server) Wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Block until a signal is received. Ctrl+C
}

func (s *Server) Run() (err error) {
	// Create a listener
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	s.stopped = false

	go func() {
		defer func() {
			s.ch <- err // Send a signal
		}()

		for { // run until listener closed
			if s.stopped {
				break
			}

			// Wait for new connection
			conn, err := s.listener.Accept()
			if err != nil {
				//log.Debug("Accept error: %v", err)
				continue
			}

			log.Debug("Accepted %v", conn.RemoteAddr())

			// Serve the connection in another goroutine
			go func() {
				h := s.handler(conn)
				r := rpc.NewServer()
				r.Register(h)

				r.ServeConn(conn) // block until disconnected

				d, ok := h.(common.OnDisconnected)
				if ok {
					d.OnDisconnected()
				}
			}()
		}
	}()

	return nil
}
