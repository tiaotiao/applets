package server

import (
	"common"
	"common/log"
	"errors"
	"io/ioutil"
	"path/filepath"
	"sync"
)

// FileManager is responsible for managing local files of peer
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

// Obtain loads and returns the file content
func (m *FileManager) Obtain(fileName string, content *[]byte) error {
	m.lock.RLock()
	f, ok := m.files[fileName]
	m.lock.RUnlock()

	if !ok {
		log.Debug("File not found in local map %v", fileName)
		return errors.New("file not found")
	}

	bytes, err := ioutil.ReadFile(f.Path) // Load file from disk
	if err != nil {
		log.Debug("Read file error=%v %v", err, f.Path)
		return err
	}

	*content = bytes
	log.Debug("Upload file ok. '%v' %v", fileName, len(bytes))

	return nil
}

// AddFolder is used for adding a whole folder of files to peer
// the whole folder will be shared and registered to Central Server
func (m *FileManager) AddFolder(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		m.AddFile(filepath.Join(folderPath, file.Name())) // Add one by one
	}
	return nil
}

// AddFile is used for adding a local file to peer
// this file will be shared and registered to Central Server
func (m *FileManager) AddFile(path string) error {
	f, err := common.GetFileInfo(path)
	if err != nil {
		return err
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	_, exist := m.files[f.Name]
	if exist { // Already added
		return nil
	}

	m.files[f.Name] = f

	// Nodify Central Server
	_, err = m.server.centralServer.Registry(&f.FileInfo)
	if err != nil {
		return err
	}

	return nil
}
