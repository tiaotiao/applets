package main

import (
	"common/log"
	"common/rpc"
	"net"
)

type LockServer struct {
	rpc     *rpc.Server
	lockMgr *LockManager
}

func NewLockServer(port int) *LockServer {
	s := &LockServer{}
	s.lockMgr = NewLockManager()

	s.rpc = rpc.NewServer(port, func(c net.Conn) interface{} {
		return NewHandler(s.lockMgr)
	})

	return s
}

func (s *LockServer) Run() error {
	err := s.rpc.Run()
	if err != nil {
		log.Error("Lock Server run failed %v", err.Error())
		return err
	}
	log.Info("Lock Server start OK.")

	s.rpc.Wait() // block

	return nil
}

func (s *LockServer) Stop() {
	s.rpc.Stop()
	log.Info("Lock Server stoped.")
}
