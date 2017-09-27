package main

import (
	"centralserver/proxy"
	"common/commander"
	"common/log"
	"errors"
	"peer/client"
	"peer/server"
)

type Peer struct {
	Server *server.Server
	Client *client.Client
	Proxy  *proxy.Proxy
	Cmd    *commander.Commander

	PeerId     string
	Port       int
	Dir        string
	ServerAddr string
}

func NewPeer(peerId string, port int, folderDir string, centralServerAddr string) *Peer {
	p := &Peer{}
	p.Proxy = proxy.NewProxy(centralServerAddr)
	p.PeerId = peerId
	p.Port = port
	p.Dir = folderDir
	p.ServerAddr = centralServerAddr

	// register cmds
	p.Cmd = commander.NewCommander()
	p.Cmd.Register("add", p.cmdAddFile, "[filepath] Add local file and register to central server")
	p.Cmd.Register("search", p.cmdSearch, "[filename] Search file from central server")
	p.Cmd.Register("obtain", p.cmdObtain, "[filename] Search and download a file from another peer")
	return p
}

func (p *Peer) cmdAddFile(args ...string) error {
	for _, filepath := range args {
		err := p.Server.FileMgr.AddFile(filepath)
		if err != nil {
			return err
		}
		log.Info("Add file ok: %s", filepath)
	}
	return nil
}

func (p *Peer) cmdSearch(args ...string) error {
	for _, filename := range args {
		results, err := p.Proxy.Search(filename)
		if err != nil {
			return err
		}
		if results == nil {
			return errors.New("file not found " + filename)
		} else {
			log.Info("Search %s: %v %v", filename, results.FileInfo, results.Peers)
		}
	}
	return nil
}

func (p *Peer) cmdObtain(args ...string) error {
	for _, filename := range args {
		err := p.Client.Obtain(filename)
		if err != nil {
			return err
		}
		//log.Info("[Cmd] Obtain file %v ok.", filename)
	}
	return nil
}

func (p *Peer) Run() {
	err := p.Proxy.Connect()
	if err != nil {
		log.Error("Connect to Central Server failed. %v, %v", p.ServerAddr, err.Error())
		return
	}
	p.Server = server.NewServer(p.Port, p.PeerId, p.Proxy)
	p.Client = client.NewClient(p.Proxy, p.Dir)

	err = p.Server.Run()
	if err != nil {
		log.Error("Start server failed. %v", err.Error())
		return
	}

	p.Server.FileMgr.AddFolder(p.Dir)

	log.Info("[%s] Start listening %v %v", p.PeerId, p.Port, p.Dir)

	p.Cmd.Run()

	p.Server.Stop()

	log.Info("Stopped.")
}
