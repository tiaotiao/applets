package main

import (
	"common"
	"common/log"
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

// Server is responseible for listening to a port and serve the connections
type Server struct {
	listener  net.Listener
	rpcServer *rpc.Server
	addr      string

	indexing *Indexing
	stoped   bool
	ch       chan error
}

func NewServer() *Server {
	s := &Server{}
	s.addr = ":" + fmt.Sprintf("%d", common.CentralServerPort)
	s.rpcServer = rpc.NewServer()
	s.indexing = NewIndexing()
	s.stoped = false
	s.ch = make(chan error, 1)
	return s
}

func (s *Server) Stop() error {
	if s.stoped {
		return errors.New("already stoped")
	}

	s.stoped = true
	s.listener.Close()

	err := <-s.ch
	return err
}

func (s *Server) Run() (err error) {
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		log.Error("Listen failed [%v], err=%v", s.addr, err)
		return err
	}
	log.Info("Start Listening %v", s.addr)

	s.stoped = false

	go func() {
		defer func() {
			s.ch <- err
		}()

		for {
			if s.stoped {
				break
			}

			// Received a new connection
			conn, err := s.listener.Accept()
			if err != nil {
				log.Debug("Accept %v", err)
				continue
			}

			log.Debug("Accepted %v", conn.RemoteAddr().String())

			h := NewHandler(conn, s.indexing)

			rpcServer := rpc.NewServer()
			rpcServer.Register(h)

			// Serve the connection in other goroutine (lightweight thread).
			go func() {
				rpcServer.ServeConn(conn) // block until disconnected
				h.onDisconnected()        // clean up all registered files
			}()
		}

		log.Info("Stopped.")
	}()

	return nil
}
