package main

import (
	"net"

	"../common"
	"../common/log"
)

type Handler struct {
	conn net.Conn
	idx  *Indexing
}

func NewHandler(conn net.Conn, idx *Indexing) *Handler {
	h := &Handler{}
	h.conn = conn
	h.idx = idx
	return h
}

///////////////////////////////////////////////
// rpc interfaces

func (h *Handler) Registry(args *common.RegistryArgs, ok *bool) error {
	err := h.idx.Registry(h.conn, args, ok)
	if err != nil {
		log.Error("[Registry] args=%v, err=%v", err)
		return err
	}
	log.Info("[Registry] args=%v, ok=&%v")
	return nil
}

func (h *Handler) Search(fileName string, results *common.SearchResults) error {
	err := h.idx.Search(fileName, results)
	if err != nil {
		log.Error("[Search] name=%v, err=%v", fileName, err)
		return err
	}
	log.Info("[Search] name=%v, result=%v", fileName, results)
	return nil
}
