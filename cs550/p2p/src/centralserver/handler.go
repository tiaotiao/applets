package main

import (
	"net"

	"common"
	"common/log"
)

type Handler struct {
	conn   net.Conn
	idx    *Indexing
	peerId string
}

func NewHandler(conn net.Conn, idx *Indexing) *Handler {
	h := &Handler{}
	h.conn = conn
	h.idx = idx
	h.peerId = ""
	return h
}

///////////////////////////////////////////////
// RPC interfaces

// Registry is the interface for RPC
func (h *Handler) Registry(args *common.RegistryArgs, ok *bool) error {
	// Get remote address
	ip := h.conn.RemoteAddr().(*net.TCPAddr).IP
	addr := ip.String()
	if ip.To4() == nil {
		addr = "[" + ip.String() + "]"
	}

	h.peerId = args.PeerId

	*ok = h.idx.Registry(addr, args) // Register to the indexing module

	log.Debug("Registry '%v', peerId=%v, size=%v, md5=%v", args.Name, args.PeerId, args.Size, args.Md5)
	return nil
}

var NotFoundResult = &common.SearchResults{Exist: false}

// Search is the interface for RPC
func (h *Handler) Search(fileName string, results *common.SearchResults) error {
	r := h.idx.Search(fileName) // Search from indexing module
	if r == nil {
		*results = *NotFoundResult
		log.Debug("Search '%v', Not found", fileName)
		return nil
	}
	*results = *r
	log.Debug("Search '%v', result=%v", fileName, results)
	return nil
}

/////////////////////////////////////////////////////////////////////////
// unexported functions

func (h *Handler) onDisconnected() {
	h.idx.RemoveAll(h.peerId)
	log.Debug("Disconnected. peerId=%v", h.peerId)
}
