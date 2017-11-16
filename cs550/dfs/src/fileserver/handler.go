package fileserver

type ClientHandler struct {
	// TODO
}

func NewClientHandler() *ClientHandler {
	// TODO
	return nil
}

///////////////////////////////////////////////
// RPC interfaces

func (h *ClientHandler) Open(path string, result interface{}) error {
	// TODO
	return nil
}

func (h *ClientHandler) Read(offset int, n int) error {
	// TODO
	return nil
}

func (h *ClientHandler) Write(data string) error {
	// TODO
	return nil
}

func (h *Clienthandler) Close() error {
	// TODO
	return nil
}

