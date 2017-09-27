package server

import (
	"common"
	"common/log"
	"errors"
	"io/ioutil"
	"path/filepath"
	"sync"
)

///////////////////////////////////////////////

type FileManager struct {
	files map[string]*common.LocalFileInfo
	lock  sync.RWMutex

	server *Server
}

func NewFileManager(s *Server) *FileManager {
	m := &FileManager{}
	m.files = make(map[string]*common.LocalFileInfo)
	m.server = s
	return m
}

func (m *FileManager) Obtain(fileName string, content *[]byte) error {
	m.lock.RLock()
	f, ok := m.files[fileName]
	m.lock.RUnlock()

	if !ok {
		log.Debug("File not found in local map %v, %v", fileName, m.files)
		return errors.New("file not found")
	}

	bytes, err := ioutil.ReadFile(f.Path)
	if err != nil {
		log.Debug("Read file error=%v %v", err, f.Path)
		return err
	}

	*content = bytes
	log.Debug("Return file ok. %v", len(bytes))

	return nil
}

func (m *FileManager) AddFolder(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		m.AddFile(filepath.Join(folderPath, file.Name()))
	}
	return nil
}

func (m *FileManager) AddFile(path string) error {
	f, err := common.GetFileInfo(path)
	if err != nil {
		return err
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	_, exist := m.files[f.Name]
	if exist {
		// TODO log warning
		return nil
	}

	m.files[f.Name] = f

	// nodify central server
	ok, err := m.server.centralServer.Registry(m.server.peerId, m.server.port, &f.FileInfo)
	if err != nil {
		return err
	}

	if !ok {
		// TODO log warning
	}

	return nil
}

// func (m *FileManager) removeFile(fileName string) error {
// 	return nil
// }
