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
// rpc interfaces

func (h *Handler) Obtain(fileName string, content *[]byte) error {
	return h.fileMgr.Obtain(fileName, content)
}
