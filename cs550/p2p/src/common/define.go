package common

var CentralServerPort int = 8099

type FileInfo struct {
	Name string
	Size int64
	Md5  string
}

func (f *FileInfo) String() string {
	return f.Name + "(" + string(f.Size) + "," + f.Md5 + ")"
}

type PeerInfo struct {
	PeerId  string
	Address string
}

func (p *PeerInfo) String() string {
	return p.PeerId + ":" + p.Address
}

///////////////////////////////////

type RegistryArgs struct {
	FileInfo
	PeerId string
}

type SearchResults struct {
	FileInfo
	Peers []PeerInfo
}

func (r *SearchResults) String() string {
	s := r.FileInfo.String()
	s += " peers(" + string(len(r.Peers)) + ")["
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
