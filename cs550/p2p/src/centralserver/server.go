package main

import "net"
import "net/rpc"
import "../common"
import "../common/log"

type Server struct {
	listener  net.Listener
	rpcServer *rpc.Server
	addr      string

	indexing *Indexing
}

func NewServer() *Server {
	s := &Server{}
	s.addr = ":" + string(common.CentralServerPort)
	s.rpcServer = rpc.NewServer()
	s.indexing = NewIndexing()

	return s
}

func (s *Server) Run() (err error) {
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		log.Error("[Server] Listen failed %s, err=%v", s.addr, err)
		return err
	}
	log.Info("[Sevrer] Start Listening %v", s.addr)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Warning("[Server] Accept error: %v", err)
		}
		log.Debug("[Server] Accepted %v", conn.RemoteAddr().String())

		h := NewHandler(conn, s.indexing)

		rpcServer := rpc.NewServer()
		rpcServer.Register(h)

		go rpcServer.ServeConn(conn)
	}

	log.Info("[Server] Stoped.")
	return nil
}
