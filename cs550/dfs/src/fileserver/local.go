package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// LocalFiles module manages local files only.
// It doesn't deal with golbal lock and replication
type LocalFiles struct {
	root string
	// TODO cache file handlers for future use.
}

func NewLocalFiles(root string) *LocalFiles {
	l := &LocalFiles{}
	l.root = root
	return nil
}

func (l *LocalFiles) Read(file string, offset int, n int) ([]byte, error) {
	// TODO optimize
	path := filepath.Join(l.root, file)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *LocalFiles) GetInfo(file string) (int, error) {
	path := filepath.Join(l.root, file)
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	s, err := f.Stat()
	if err != nil {
		return 0, err
	}
	size := int(s.Size())
	return size, nil
}

func (l *LocalFiles) Append(file string, data []byte) error {
	path := filepath.Join(l.root, file)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	n, err := f.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		println("n!=len(data)", n, len(data))
	}

	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
