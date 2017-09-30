package main

import (
	"common/log"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func readParams() (server string, dir string, port int, peerId string, debug bool) {
	flag.StringVar(&server, "server", "localhost", "Address of central server. Default value is localhost.")
	flag.StringVar(&dir, "dir", "./", "Shared directory. Default is current dir.")
	flag.IntVar(&port, "port", 0, "Optional. Listening port number. Default port is chosen randomly.")
	flag.StringVar(&peerId, "id", "", "Optional. A random ID will be given if not specified.")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if port == 0 {
		port = 8100 + rand.Intn(900) // random port
	}

	if peerId == "" {
		peerId = fmt.Sprintf("peer-%d-%s", port, randString(4)) // random peer id
	}
	return
}

func randString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	server, dir, port, peerId, debug := readParams()

	log.LevelDebug = debug
	log.ModuleName = "Peer"

	p := NewPeer(peerId, port, dir, server)

	p.Run()
}
