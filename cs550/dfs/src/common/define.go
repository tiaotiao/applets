package common

const LOCKSERVER_PORT = 4090

type OnDisconnected interface {
	OnDisconnected()
}

type ParamUpdate struct {
	FileID string
	Data   []byte
}

///////////////////////////////////////////////////

type OpenResult struct {
	FileID string
	Size   int
}

type ReadParam struct {
	FileID string
	Offset int
	Count  int
}

type WriteParam struct {
	FileID string
	Data   []byte
}
