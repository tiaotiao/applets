package main

import (
	"errors"
	"time"

	"github.com/tiaotiao/mapstruct"
)

type ClientHandler struct {
	fileserver  *FileServer
	openedFiles map[string]int
}

func NewClientHandler(fileserver *FileServer) *ClientHandler {
	h := &ClientHandler{}
	h.fileserver = fileserver
	h.openedFiles = make(map[string]int)
	return h
}

///////////////////////////////////////////////
// RPC interfaces

func (h *ClientHandler) Open(path string, results *map[string]interface{}) error {
	fileID, size, err := h.fileserver.Open(path)
	if err != nil {
		return err
	}

	h.openedFiles[fileID] = time.Now().Nanosecond()

	r := make(map[string]interface{})
	r["fileid"] = fileID
	r["size"] = size
	*results = r
	return nil
}

func (h *ClientHandler) Read(params map[string]interface{}, data *[]byte) error {
	s := struct {
		FileID string `map:"fileid,required"`
		Offset int    `map:"offset,required"`
		Count  int    `map:"n,required"`
	}{}

	err := mapstruct.Map2Struct(params, &s)
	if err != nil {
		return err
	}

	_, ok := h.openedFiles[s.FileID]
	if !ok {
		return errors.New("not opend " + s.FileID)
	}

	d, err := h.fileserver.Read(s.FileID, s.Offset, s.Count)
	if err != nil {
		return err
	}

	*data = d
	return nil
}

func (h *ClientHandler) Write(params map[string]interface{}, n *int) error {
	s := struct {
		FileID string `map:"fileid,required"`
		Data   []byte `map:"data,required"`
	}{}

	err := mapstruct.Map2Struct(params, &s)
	if err != nil {
		return err
	}

	_, ok := h.openedFiles[s.FileID]
	if !ok {
		return errors.New("not opend " + s.FileID)
	}

	err = h.fileserver.Write(s.FileID, s.Data)
	if err != nil {
		return err
	}
	*n = len(s.Data) // TODO
	return nil
}

func (h *ClientHandler) Close(fileID string, ok *bool) error {
	_, *ok = h.openedFiles[fileID]
	if *ok {
		delete(h.openedFiles, fileID)
	}
	return nil
}
