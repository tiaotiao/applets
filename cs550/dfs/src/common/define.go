package common

const LOCKSERVER_PORT = 4090

type OnDisconnected interface {
	OnDisconnected()
}

type ParamUpdate struct {
	FileID string
	Data   []byte
}
