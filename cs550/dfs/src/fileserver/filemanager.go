package fileserver

import "lockserver/sdk" 

type FileManager struct {
	local *LocalFiles
	lockServer *sdk.LockServerClient
	// TODO
}

func NewFileManager() *FileManager {
	// TODO
	return nil
}

func (m *FileManager) Open(path string) (token string, err error) {
	// TODO
	// convert path to token
	return
}

func (m *FileManager) Read(token string, offset int, n int) ([]byte, error) {
	// TODO
	// m.lockServer.Require(path, sdk.PermRead, 0)
	// m.local.Read(token, offset, n)
	// m.lockServer.Release(path)
	return nil
}

func (m *FileManager) Write(token string, data []byte) error {
	// TODO
	// m.lockServer.Require(path, sdk.PermWrite, 0)
	// m.local.Append(token, data)
	// update file of the next two fileservers
	// remember the token to be released
	// wait for the update request came back, then m.lockServer.Release(path)
}

func (m *FileManager) Update(token string, data []byte) error {
	// TODO 
	// If get back, release the write lock
	// else
	// m.local.Append(token, data)
	// update file of the next two fileservers
}