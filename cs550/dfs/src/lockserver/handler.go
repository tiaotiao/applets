package main

import (
	"common/log"
	"fmt"
	"sync"
)

var g_connID int = 0

type Handler struct {
	connID  string
	lockMgr *LockManager

	// keys  map[string]bool
	mutex sync.Mutex
}

func NewHandler(lockMgr *LockManager) *Handler {
	h := &Handler{}
	h.lockMgr = lockMgr
	// h.keys = make(map[string]bool)
	// h.connID = common.RandString(16)

	h.mutex.Lock()
	g_connID += 1
	id := g_connID
	h.mutex.Unlock()

	h.connID = fmt.Sprintf("%d", id)
	return h
}

// OnDisconnected implements the interface
func (h *Handler) OnDisconnected() {
	h.lockMgr.RelaseUser(h.connID)
	log.Info("Disconnected conn=%v", h.connID)
}

////////////////////////////////////////////////////////////////
// handlers for RPC

func (h *Handler) AcquireRead(key string, ok *bool) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	*ok = h.lockMgr.Acquire(h.connID, key, false)
	log.Info("Acquire [Read] conn=%v, key=(%v), %v", h.connID, key, *ok)
	log.Info("Lock status %v", h.lockMgr.String())

	return nil
}

func (h *Handler) AcquireWrite(key string, ok *bool) error {
	*ok = h.lockMgr.Acquire(h.connID, key, true)
	log.Info("Acquire [Write] conn=%v, key=(%v), %v", h.connID, key, *ok)
	log.Info("Lock status %v", h.lockMgr.String())
	return nil
}

func (h *Handler) Release(key string, ok *bool) error {
	h.lockMgr.Relase(h.connID, key)
	*ok = true
	log.Info("Release conn=%v, key=(%v), %v", h.connID, key, *ok)
	log.Info("Lock status %v", h.lockMgr.String())
	return nil
}
