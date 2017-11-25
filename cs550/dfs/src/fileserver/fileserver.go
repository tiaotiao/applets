package main

import (
	"common/log"
	"common/rpc"
	"crypto/md5"
	"errors"
	fssdk "fileserver/sdk"
	"fmt"
	"lockserver/sdk"
	"net"
)

type FileServer struct {
	local *LocalFiles
	locks *sdk.LockClient

	cliRpc *rpc.Server // external interface for clients
	svrRpc *rpc.Server // internal interface for other fileservers

	protocol *fssdk.FileProtocol
}

func NewFileServer(externalPort int, internalPort int, lockServerAddr string, fileServerAddrs []string) *FileServer {
	s := &FileServer{}

	s.cliRpc = rpc.NewServer(externalPort, func(c net.Conn) interface{} {
		return NewClientHandler(s)
	})

	s.svrRpc = rpc.NewServer(internalPort, func(c net.Conn) interface{} {
		return NewFileServerProtocol(s)
	})

	s.locks = sdk.NewLockClient(lockServerAddr)

	s.protocol = fssdk.NewFileProtocol(fileServerAddrs)

	return s
}

///////////////////////////////////////////////////////////////////////

func (s *FileServer) Run() (err error) {
	// NOTICE: The order of following operations does matter
	err = s.svrRpc.Run() // run server for fileservers
	if err != nil {
		s.cliRpc.Stop()
		return err
	}
	log.Info("Server rpc started.")

	err = s.locks.Connect() // connect to lock server
	if err != nil {
		return err
	}
	log.Info("Lock Server connected.")

	err = s.protocol.Connect() // connect to other file servers. block until all connected or failed
	if err != nil {
		return err
	}
	log.Info("File Servers connected.")

	err = s.cliRpc.Run() // run server for clients
	if err != nil {
		return err
	}
	log.Info("Client rpc started.")

	log.Info("Running OK.")
	s.cliRpc.Wait() // block

	return nil
}

func (s *FileServer) Stop() {
	s.cliRpc.Stop()
	s.svrRpc.Stop()
	s.locks.Close()
	log.Info("Server Stoped.")
}

///////////////////////////////////////////////////////////////////////////////////

func (s *FileServer) Open(path string) (fileID string, size int, err error) {
	// convert path to fileID and file info
	h := md5.New()
	fileID = fmt.Sprintf("%x", h.Sum([]byte(path)))

	size, err = s.local.GetInfo(fileID)
	if err != nil {
		size = 0
	}
	return fileID, size, nil
}

func (s *FileServer) Read(fileID string, offset int, n int) ([]byte, error) {
	ok, err := s.locks.RequireRead(fileID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("Require lock failed")
	}

	// read local file
	data, err := s.local.Read(fileID, offset, n)
	if err != nil {
		return nil, err
	}

	ok, err = s.locks.Release(fileID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *FileServer) Write(fileID string, data []byte) error {
	ok, err := s.locks.RequireWrite(fileID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Require lock failed")
	}

	// write local file
	err = s.local.Append(fileID, data)
	if err != nil {
		return err
	}

	// update other servers in background
	go func() {
		// send update requests to all other servers
		err = s.protocol.Update(fileID, data)
		if err != nil {
			log.Error("Update protocol failed, update %v %v, err=%v", fileID, string(data), err)
		}

		ok, err = s.locks.Release(fileID)
		if err != nil {
			log.Error("Release lock failed, update %v %v, err=%v", fileID, string(data), err)
		}
	}()

	// return immediately
	return nil
}

// Update deal with the protocol from other server
// File been modified by other server, just update the local file here
func (s *FileServer) Update(fileID string, data []byte) error {
	return s.local.Append(fileID, data)
}
