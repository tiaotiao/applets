package main

import (
	"flag"
)

func main() {
	var addr string
	flag.StringVar(&addr, "server", "", "address of file server")
	flag.Parse()

	c := NewClient(addr)
	err := c.Run()
	if err != nil {
		println("Run error", err.Error())
	}
}
