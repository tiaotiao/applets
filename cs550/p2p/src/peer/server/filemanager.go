package server

import (
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"../../common"
)

type LocalFileInfo struct {
	common.FileInfo
	Path string
}

///////////////////////////////////////////////

type FileManager struct {
	files map[string]*LocalFileInfo
	lock  sync.RWMutex

	server *Server
}

func NewFileManager(s *Server) *FileManager {
	m := &FileManager{}
	m.files = make(map[string]*LocalFileInfo)
	m.server = s
	return m
}

///////////////////////////////////////////////
// rpc interfaces

func (m *FileManager) Obtain(fileName string, content *[]byte) error {
	m.lock.RLock()
	f, ok := m.files[fileName]
	m.lock.RUnlock()

	if !ok {
		return errors.New("file not found")
	}

	bytes, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}

	*content = bytes

	return nil
}

///////////////////////////////////////////////
// unexproted functions

func (m *FileManager) addFolder(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		m.addFile(filepath.Join(folderPath, file.Name()))
	}
	return nil
}

func (m *FileManager) addFile(path string) error {
	f, err := getFileInfo(path)
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
	ok, err := m.server.centralServer.Registry(m.server.peerId, &f.FileInfo)
	if err != nil {
		return err
	}

	if !ok {
		// TODO log warning
	}

	return nil
}

func getFileInfo(path string) (*LocalFileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	s, err := file.Stat()
	if err != nil {
		return nil, err
	}

	f := &LocalFileInfo{}
	f.Size = s.Size()
	f.Name = s.Name()

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return nil, err
	}
	f.Md5 = string(h.Sum(nil))

	return f, nil
}

// func (m *FileManager) removeFile(fileName string) error {
// 	return nil
// }
