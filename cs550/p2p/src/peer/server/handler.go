package server

type Handler struct {
	fileMgr *FileManager
}

func NewHandler(fileMgr *FileManager) *Handler {
	h := &Handler{}
	h.fileMgr = fileMgr
	return h
}

///////////////////////////////////////////////
// RPC interfaces

// Obtain provide the RPC interface for other peer
func (h *Handler) Obtain(fileName string, content *[]byte) error {
	return h.fileMgr.Obtain(fileName, content)
}

// Echo provids a way to test the connection
func (h *Handler) Echo(hello string, echo *string) error {
	*echo = hello
	return nil
}
