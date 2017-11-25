package main

import (
	"common"
)

type FileServerProtocol struct {
	fileserver *FileServer
}

func NewFileServerProtocol(fileserver *FileServer) *FileServerProtocol {
	p := &FileServerProtocol{}
	p.fileserver = fileserver
	return p
}

////////////////////////////////////////////////////////////////////////////////
// Deal with file server protocol

func (p *FileServerProtocol) Update(param common.ParamUpdate, ok *bool) error {
	err := p.fileserver.Update(param.FileID, param.Data)
	if err != nil {
		*ok = false
		return err
	}
	*ok = true
	return nil
}
