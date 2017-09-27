package common

import "fmt"

const CentralServerPort int = 8099

type FileInfo struct {
	Name string
	Size int64
	Md5  string
}

func (f *FileInfo) String() string {
	return f.Name + " (size " + fmt.Sprintf("%d", f.Size) + ", md5 " + f.Md5 + ")"
}

type LocalFileInfo struct {
	FileInfo
	Path string
}

type PeerInfo struct {
	PeerId  string
	Address string
	Port    int
}

func (p *PeerInfo) String() string {
	return fmt.Sprintf("%v,%v:%v", p.PeerId, p.Address, p.Port)
}

///////////////////////////////////

type RegistryArgs struct {
	FileInfo
	PeerId string
	Port   int
}

type SearchResults struct {
	Exist bool
	FileInfo
	Peers []PeerInfo
}

func (r *SearchResults) String() string {
	if !r.Exist {
		return "not found"
	}
	s := r.FileInfo.String()
	s += " peers(" + fmt.Sprintf("%d", len(r.Peers)) + ")["
	for i, p := range r.Peers {
		if i != 0 {
			s += ","
		}
		s += p.String()
	}
	s += "]"
	return s
}

///////////////////////////////////
