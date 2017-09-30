package server

import (
	"centralserver/proxy"
	"fmt"
	"net"
	"net/rpc"

	"common/log"
)

// Server is responseible for listening to a port and serve the connections
type Server struct {
	listener  net.Listener
	rpcServer *rpc.Server
	port      int
	addr      string

	FileMgr       *FileManager
	centralServer *proxy.Proxy
	peerId        string
	handler       *Handler

	stopped bool
	ch      chan error
}

func NewServer(port int, peerId string, p *proxy.Proxy) *Server {
	s := &Server{}
	s.centralServer = p
	s.peerId = peerId
	s.port = port
	s.addr = ":" + fmt.Sprintf("%d", port)
	s.rpcServer = rpc.NewServer()
	s.FileMgr = NewFileManager(s)
	s.handler = NewHandler(s.FileMgr)

	s.rpcServer.Register(s.handler)
	s.stopped = false
	s.ch = make(chan error)
	return s
}

func (s *Server) Stop() {
	s.stopped = true
	s.listener.Close()
	<-s.ch
}

func (s *Server) Run() (err error) {
	s.listener, err = net.Listen("tcp", s.addr) // Listen a port
	if err != nil {
		return err
	}

	s.stopped = false

	go func() {
		defer func() {
			s.ch <- err
		}()

		for {
			if s.stopped {
				break
			}
			conn, err := s.listener.Accept() // New connection
			if err != nil {
				log.Debug("Accept error: %v", err)
				continue
			}

			log.Debug("Accepted %v", conn.RemoteAddr().String())

			go s.rpcServer.ServeConn(conn) // Serve the connection in another goroutine (lightweight thread)
		}
	}()

	return nil
}
