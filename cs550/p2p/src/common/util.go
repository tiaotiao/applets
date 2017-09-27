package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func GetFileInfo(path string) (*LocalFileInfo, error) {
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
	f.Path = path

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return nil, err
	}
	f.Md5 = fmt.Sprintf("%X", h.Sum(nil))

	return f, nil
}
