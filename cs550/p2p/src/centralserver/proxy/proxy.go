package proxy

import (
	"common"
	"math/rand"
	"net/rpc"
)

// Proxy is a client interface to communicate with Central Server
type Proxy struct {
	rpcs      []*rpc.Client
	addresses []string
	peerId    string
	peerPort  int
}

func NewProxy(servers []string, peerId string, peerPort int) *Proxy {
	if len(servers) == 0 {
		return nil
	}
	p := &Proxy{}
	p.addresses = servers
	p.peerId = peerId
	p.peerPort = peerPort
	return p
}

// Connect to Central Server
func (p *Proxy) Connect() (err error) {
	p.rpcs = make([]*rpc.Client, len(p.addresses))

	for i, addr := range p.addresses {
		rpc, err := rpc.Dial("tcp", addr)
		if err != nil {
			return err
		}
		p.rpcs[i] = rpc
	}
	return nil
}

// Registry a new file to Central Server
func (p *Proxy) Registry(f *common.FileInfo) (ok bool, err error) {
	args := &common.RegistryArgs{
		PeerId:   p.peerId,
		Port:     p.peerPort,
		FileInfo: *f,
	}

	idx := rand.Intn(len(p.rpcs))
	rpc := p.rpcs[idx]

	err = rpc.Call("Handler.Registry", args, &ok)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// Search a file by name from Central Server
func (p *Proxy) Search(fileName string) (*common.SearchResults, error) {
	var results common.SearchResults

	offset := rand.Intn(len(p.rpcs))

	for k := 0; k < len(p.rpcs); k++ {
		i := (offset + k) % len(p.rpcs)
		rpc := p.rpcs[i]
		err := rpc.Call("Handler.Search", fileName, &results)
		if err != nil {
			return nil, err
		}
		if results.Exist {
			break
		}
	}

	return &results, nil
}

func (p *Proxy) ListAll() (map[string]*common.SearchResults, error) {
	var results = make(map[string]*common.SearchResults)

	for i := 0; i < len(p.rpcs); i++ {
		var r map[string]*common.SearchResults
		rpc := p.rpcs[i]

		err := rpc.Call("Handler.ListAll", "", &r)
		if err != nil {
			return nil, err
		}

		for name, f := range r {
			file, ok := results[name]
			if !ok {
				results[name] = f
			} else {
				file.Peers = append(file.Peers, f.Peers...)
			}
		}
	}

	return results, nil
}
