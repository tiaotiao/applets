package main

import "sync"
import "errors"
import "../common"
import "net"

type Indexing struct {
	files map[string]*common.SearchResults
	lock  sync.RWMutex
}

func NewIndexing() *Indexing {
	i := &Indexing{}
	i.files = make(map[string]*common.SearchResults)
	return i
}

func (i *Indexing) Registry(conn net.Conn, args *common.RegistryArgs, ok *bool) error {
	if args == nil {
		return errors.New("args is nil")
	}

	i.lock.Lock()

	f, exist := i.files[args.Name]

	if !exist {
		f = &common.SearchResults{}
		f.Name = args.Name
		f.Size = args.Size

		i.files[args.Name] = f
	}

	*ok = true
	for _, p := range f.Peers {
		if p.PeerId == args.PeerId {
			*ok = false
			break
		}
	}
	if *ok {
		peers := make([]common.PeerInfo, len(f.Peers)+1)
		copy(peers, f.Peers)

		p := common.PeerInfo{
			PeerId:  args.PeerId,
			Address: conn.RemoteAddr().String(),
		}
		f.Peers = append(peers, p)
	}

	i.lock.Unlock()

	return nil
}

func (i *Indexing) Search(fileName string, results *common.SearchResults) error {
	i.lock.RLock()
	defer i.lock.RUnlock()

	f, exist := i.files[fileName]
	if !exist {
		return nil
	}
	results = f
	return nil
}
