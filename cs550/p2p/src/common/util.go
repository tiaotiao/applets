package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
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

func RandString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
