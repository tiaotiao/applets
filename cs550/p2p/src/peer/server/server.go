package server

import "net"
import "net/rpc"
import "../../centralserver/proxy"

type Server struct {
	listener  net.Listener
	rpcServer *rpc.Server
	addr      string

	fileMgr       *FileManager
	centralServer *proxy.Proxy
	peerId        string

	stopped bool
}

func NewServer(port int, peerId string, p *proxy.Proxy) *Server {
	s := &Server{}
	s.addr = ":" + string(port)
	s.rpcServer = rpc.NewServer()
	s.fileMgr = NewFileManager(s)

	s.rpcServer.Register(s.fileMgr)
	s.stopped = true
	return s
}

func (s *Server) Stop() {
	s.stopped = true
	s.listener.Close()
}

func (s *Server) Run() (err error) {
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.stopped = false
	for {
		if s.stopped {
			break
		}
		conn, err := s.listener.Accept()
		if err != nil {
			// TODO log error
		}

		go s.rpcServer.ServeConn(conn)
	}
	return nil
}
