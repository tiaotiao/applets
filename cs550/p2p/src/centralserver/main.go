package main

import (
	"common/log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	log.LevelDebug = cfg.Debug
	log.ModuleName = "Server"

	s := NewServer(cfg.Port)

	err = s.Run() // run server
	if err != nil {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Block until a signal is received. Ctrl+C
	s.Stop()

	return
}
