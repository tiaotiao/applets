package main

import (
	"common"
	"sync"
)

type Indexing struct {
	files map[string]*common.SearchResults       // filename -> fileInfo and list of peers
	peers map[string]map[string]*common.FileInfo // peerId -> map of files
	lock  sync.RWMutex
}

func NewIndexing() *Indexing {
	i := &Indexing{}
	i.files = make(map[string]*common.SearchResults)
	i.peers = make(map[string]map[string]*common.FileInfo)
	return i
}

// Registry add a file info to the index
func (i *Indexing) Registry(addr string, args *common.RegistryArgs) bool {
	if args == nil {
		return false
	}
	i.lock.Lock()

	// Check if fileName already exist
	f, exist := i.files[args.Name]

	if !exist {
		f = &common.SearchResults{}
		f.Exist = true
		f.FileInfo = args.FileInfo

		i.files[args.Name] = f
	}

	// Check if the file has already registered by this peer
	ok := true
	for i, p := range f.Peers {
		if p.PeerId == args.PeerId {
			ok = false
			// Update peer info
			f.Peers[i].Address = addr
			f.Peers[i].Port = args.Port
			break
		}
	}

	if ok {
		p := common.PeerInfo{
			PeerId:  args.PeerId,
			Address: addr,
			Port:    args.Port,
		}
		// Append the peer
		f.Peers = append(f.Peers, p)

		if i.peers[p.PeerId] == nil {
			i.peers[p.PeerId] = make(map[string]*common.FileInfo)
		}
		i.peers[p.PeerId][f.Name] = &f.FileInfo
	}

	i.lock.Unlock()

	return ok
}

// Search file info by name
func (i *Indexing) Search(fileName string) *common.SearchResults {
	i.lock.RLock()
	defer i.lock.RUnlock()

	f, exist := i.files[fileName]
	if !exist {
		return nil
	}

	return f
}

func (i *Indexing) ListAll() map[string]*common.SearchResults {
	i.lock.RLock()
	defer i.lock.RUnlock()

	return i.files
}

// Remove a file from index. Unused for now.
func (i *Indexing) Remove(fileName string, peerId string) bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	files := i.peers[peerId]
	if files != nil {
		delete(files, fileName)
	}

	r := i.files[fileName]
	if r == nil {
		return false
	}

	for j, p := range r.Peers {
		if p.PeerId == peerId {
			// remove item from array
			r.Peers[j] = r.Peers[len(r.Peers)-1]
			r.Peers = r.Peers[:len(r.Peers)-1]
			if len(r.Peers) <= 0 {
				delete(i.files, fileName)
			}
			return true
		}
	}

	return false
}

// RemoveAll registered files by peerId
func (i *Indexing) RemoveAll(peerId string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	files := i.peers[peerId]
	delete(i.peers, peerId)

	if files == nil {
		return
	}

	for name := range files {
		r := i.files[name]
		if r == nil {
			continue
		}
		for j, p := range r.Peers {
			if p.PeerId == peerId {
				// Remove item from array
				r.Peers[j] = r.Peers[len(r.Peers)-1]
				r.Peers = r.Peers[:len(r.Peers)-1]
				if len(r.Peers) <= 0 {
					delete(i.files, name)
				}
				break
			}
		}

	}
}
