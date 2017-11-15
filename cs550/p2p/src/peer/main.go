package main

import (
	"common/log"
	"fmt"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	log.LevelDebug = cfg.Debug
	log.ModuleName = "Peer"

	p := NewPeer(cfg.PeerID, cfg.Port, cfg.Dir, cfg.Servers)

	p.Run()
}
