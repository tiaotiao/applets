package sdk

import (
	"common"
	"common/log"
	"net/rpc"
	"time"
)

var MAX_RETRY = 30

type FileProtocol struct {
	addresses []string
	rpcs      []*rpc.Client
}

func NewFileProtocol(addrs []string) *FileProtocol {
	p := &FileProtocol{}
	p.addresses = addrs
	return p
}

func (p *FileProtocol) Connect() (err error) {
	for _, addr := range p.addresses {
		err = p.connect(addr, MAX_RETRY)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *FileProtocol) connect(addr string, maxRetry int) error {
	retry := 0
	for {
		rpc, err := rpc.Dial("tcp", addr)
		if err != nil {
			retry++
			if retry > MAX_RETRY {
				return err
			}
			<-time.After(time.Second) // sleep for a second
			continue
		}
		p.rpcs = append(p.rpcs, rpc)
		break
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////
// RPC protocols

func (p *FileProtocol) Update(fileID string, data []byte) error {
	param := &common.ParamUpdate{}
	param.FileID = fileID
	param.Data = data

	for i, rpc := range p.rpcs { // TODO retry
		ok, err := p.update(rpc, param)
		if err != nil || !ok {
			log.Error("Update failed %v, %v %v", p.addresses[i], ok, err)
		}
	}

	return nil
}

func (p *FileProtocol) update(rpc *rpc.Client, param *common.ParamUpdate) (bool, error) {
	var ok bool
	err := rpc.Call("FileServerProtocol.Update", param, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}
