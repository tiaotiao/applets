package main

import (
	"common/log"
	"sync"
)

type Handler struct {
	connID  string
	lockMgr *LockManager

	keys  map[string]bool
	mutex sync.Mutex
}

func NewHandler(lockMgr *LockManager) *Handler {
	h := &Handler{}
	h.lockMgr = lockMgr
	return h
}

// OnDisconnected implements the interface
func (h *Handler) OnDisconnected() {
	h.lockMgr.RelaseUser(h.connID)
	log.Info("Disconnected conn=%v", h.connID)
}

////////////////////////////////////////////////////////////////
// handlers for RPC

func (h *Handler) RequireRead(key string, ok *bool) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	*ok = h.lockMgr.Require(h.connID, key, false)
	log.Info("Require [Read] conn=%v, key=(%v), %v", h.connID, key, *ok)

	return nil
}

func (h *Handler) RequireWrite(key string, ok *bool) error {
	*ok = h.lockMgr.Require(h.connID, key, true)
	log.Info("Require [Write] conn=%v, key=(%v), %v", h.connID, key, *ok)
	return nil
}

func (h *Handler) Release(key string, ok *bool) error {
	h.lockMgr.Relase(h.connID, key)
	*ok = true
	log.Info("Release conn=%v, key=(%v), %v", h.connID, key, *ok)
	return nil
}
