package main

import (
	"common/log"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.ModuleName = "Server"
}

func main() {
	flag.Parse()

	s := NewServer()

	err := s.Run()
	if err != nil {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block until a signal is received. Ctrl+C
	s.Stop()

	return
}
