package main

import (
	"centralserver/proxy"
	"common"
	"common/commander"
	"common/log"
	"fmt"
	"math/rand"
	"peer/client"
	"peer/server"
	"time"
)

// Peer is the top level object managing all the modules
type Peer struct {
	Server *server.Server
	Client *client.Client
	Proxy  *proxy.Proxy
	Cmd    *commander.Commander

	PeerId      string
	Port        int
	Dir         string
	ServerAddrs []string
}

func NewPeer(peerId string, port int, folderDir string, servers []string) *Peer {
	rand.Seed(time.Now().UnixNano())
	if port == 0 {
		port = 8100 + rand.Intn(900) // random port
	}
	if peerId == "" {
		peerId = fmt.Sprintf("PEER-%d-%s", port, common.RandString(4)) // random peer id
	}

	p := &Peer{}
	p.Proxy = proxy.NewProxy(servers, peerId, port)
	p.PeerId = peerId
	p.Port = port
	p.Dir = folderDir
	p.ServerAddrs = servers

	// Register commands
	p.Cmd = commander.NewCommander()
	p.Cmd.Register("add", p.cmdAddFile, "[filepath] Add local file and register to central server")
	p.Cmd.Register("search", p.cmdSearch, "[filename] Search file from central server")
	p.Cmd.Register("list", p.cmdListAll, "List all file from central server")
	p.Cmd.Register("obtain", p.cmdObtain, "[filename] Search and download a file from another peer")
	p.Cmd.Register("test", p.cmdTestProformance, "[n] Test search for n times. Default n is 10000.")
	return p
}

func (p *Peer) Run() {
	// Connected to Central Server
	err := p.Proxy.Connect()
	if err != nil {
		log.Error("Connect to Central Server failed. %v, %v", p.ServerAddrs, err.Error())
		return
	}
	p.Server = server.NewServer(p.Port, p.PeerId, p.Proxy)
	p.Client = client.NewClient(p.PeerId, p.Proxy, p.Dir)

	// Start peer server
	err = p.Server.Run()
	if err != nil {
		log.Error("Start server failed. %v", err.Error())
		return
	}

	// Add local files and register them to Central Server
	p.Server.FileMgr.AddFolder(p.Dir)

	log.Info("[%s] Start listening %v %v", p.PeerId, p.Port, p.Dir)

	// Start command line, block until exit
	p.Cmd.Run()

	p.Server.Stop()

	log.Info("Stopped.")
}
