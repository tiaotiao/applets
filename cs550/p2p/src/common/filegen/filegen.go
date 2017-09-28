package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
)

func generateFile(path string, size int) error {
	b := make([]byte, size/2)

	rand.Read(b)

	h := make([]byte, size)
	hex.Encode(h, b)

	return ioutil.WriteFile(path, h, 0666)
}

func generateFiles(dir string, size int, step int, count int) error {
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%d_%d.txt", i, size)

		err := generateFile(path.Join(dir, name), size)
		if err != nil {
			return err
		}

		fmt.Printf("File generated. %v\n", name)

		size += step
	}
	return nil
}

func readParams() (dir string, size int, step int, count int) {
	flag.StringVar(&dir, "dir", "./", "Directory path where files are generated to.")
	flag.IntVar(&size, "size", 1024, "Size of file.")
	flag.IntVar(&step, "step", 0, "Step of size while increasing.")
	flag.IntVar(&count, "count", 1, "Number of files to be generated.")
	flag.Parse()
	return
}

func main() {
	dir, size, step, count := readParams()

	err := generateFiles(dir, size, step, count)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
