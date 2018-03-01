package main

import (
	"common"
	"errors"
	"time"
)

type ClientHandler struct {
	fileserver      *FileServer
	openedFiles     map[string]int
	lockedFileID    string
	lockedExclusive bool
}

func NewClientHandler(fileserver *FileServer) *ClientHandler {
	h := &ClientHandler{}
	h.fileserver = fileserver
	h.openedFiles = make(map[string]int)
	h.lockedFileID = ""
	h.lockedExclusive = false
	return h
}

///////////////////////////////////////////////
// RPC interfaces

func (h *ClientHandler) Open(path string, result *common.OpenResult) error {
	fileID, size, err := h.fileserver.Open(path)
	if err != nil {
		return err
	}

	h.openedFiles[fileID] = time.Now().Nanosecond()

	(*result).FileID = fileID
	(*result).Size = size
	return nil
}

func (h *ClientHandler) Read(param common.ReadParam, data *[]byte) error {
	_, ok := h.openedFiles[param.FileID]
	if !ok {
		return errors.New("not opend " + param.FileID)
	}

	d, err := h.fileserver.Read(param.FileID, param.Offset, param.Count)
	if err != nil {
		return err
	}

	*data = d
	return nil
}

func (h *ClientHandler) Write(param common.WriteParam, n *int) error {
	_, ok := h.openedFiles[param.FileID]
	if !ok {
		return errors.New("not opend " + param.FileID)
	}

	err := h.fileserver.Write(param.FileID, param.Data)
	if err != nil {
		return err
	}

	*n = len(param.Data) // TODO
	return nil
}

func (h *ClientHandler) AcquireRead(fileID string, ok *bool) (err error) {
	if h.lockedFileID != "" {
		if h.lockedFileID == fileID && h.lockedExclusive == false {
			*ok = true
		} else {
			*ok = false
		}
		return nil
	}

	*ok, err = h.fileserver.locks.AcquireRead(fileID)
	if err != nil {
		return err
	}
	if !*ok {
		return nil
	}

	h.lockedFileID = fileID
	h.lockedExclusive = false
	*ok = true
	return nil
}

func (h *ClientHandler) AcquireWrite(fileID string, ok *bool) (err error) {
	if h.lockedFileID != "" {
		*ok = false
		return nil
	}

	*ok, err = h.fileserver.locks.AcquireWrite(fileID)
	if err != nil {
		return err
	}
	if !*ok {
		return nil
	}

	h.lockedFileID = fileID
	h.lockedExclusive = true
	return nil
}

func (h *ClientHandler) Release(fileID string, ok *bool) (err error) {
	if h.lockedFileID == "" {
		*ok = false
		return nil
	}
	if h.lockedFileID != fileID {
		*ok = false
		return nil
	}

	h.lockedFileID = ""
	h.lockedExclusive = false

	*ok, err = h.fileserver.locks.Release(fileID)
	if err != nil {
		return err
	}
	return nil
}

func (h *ClientHandler) Close(fileID string, ok *bool) error {
	_, *ok = h.openedFiles[fileID]
	if *ok {
		delete(h.openedFiles, fileID)
	}
	return nil
}
