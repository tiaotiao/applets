package sdk

import (
	"common"
	"fmt"
	"net/rpc"
)

type LockClient struct {
	address string
	rpc     *rpc.Client
}

func NewLockClient(address string) *LockClient {
	l := LockClient{}
	l.address = fmt.Sprintf("%s:%d", address, common.LOCKSERVER_PORT)
	return &l
}

func (l *LockClient) Connect() (err error) {
	l.rpc, err = rpc.Dial("tcp", l.address)
	if err != nil {
		return err
	}
	return nil
}

func (l *LockClient) Close() error {
	return l.rpc.Close()
}

///////////////////////////////////////////////////////////////////////////

func (l *LockClient) RequireRead(fileID string) (bool, error) {
	var ok bool
	err := l.rpc.Call("Handler.RequireRead", fileID, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (l *LockClient) RequireWrite(fileID string) (bool, error) {
	var ok bool
	err := l.rpc.Call("Handler.RequireWrite", fileID, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (l *LockClient) Release(fileID string) (bool, error) {
	var ok bool
	err := l.rpc.Call("Handler.Release", fileID, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}
