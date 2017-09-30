package proxy

import (
	"common"
	"fmt"
	"net/rpc"
)

// Proxy is a client interface to communicate with Central Server
type Proxy struct {
	rpcClient *rpc.Client
	addr      string
}

func NewProxy(addr string) *Proxy {
	p := &Proxy{}
	p.addr = addr + ":" + fmt.Sprintf("%d", common.CentralServerPort)
	return p
}

// Connect to Central Server
func (p *Proxy) Connect() (err error) {
	p.rpcClient, err = rpc.Dial("tcp", p.addr)
	if err != nil {
		return err
	}
	return nil
}

// Registry a new file to Central Server
func (p *Proxy) Registry(peerId string, port int, f *common.FileInfo) (ok bool, err error) {
	args := &common.RegistryArgs{
		PeerId:   peerId,
		Port:     port,
		FileInfo: *f,
	}

	err = p.rpcClient.Call("Handler.Registry", args, &ok)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// Search a file by name from Central Server
func (p *Proxy) Search(fileName string) (*common.SearchResults, error) {
	var results common.SearchResults
	err := p.rpcClient.Call("Handler.Search", fileName, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
