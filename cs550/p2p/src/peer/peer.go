package main

import (
	"centralserver/proxy"
	"common/commander"
	"common/log"
	"errors"
	"peer/client"
	"peer/server"
	"strconv"
	"time"
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
	p.Cmd.Register("test", p.cmdTestSearch, "[filename n] Test search for n times. Default n is 1000.")
	return p
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

///////////////////////////////////////////////

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
		if !results.Exist {
			log.Info("Search '%s' not found", filename)
		} else {
			log.Info("Search '%s' size=%v, peers=%v", filename, results.Size, results.Peers)
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

func (p *Peer) cmdTestSearch(args ...string) (err error) {
	n := 1000
	if len(args) < 0 {
		return nil
	}
	filename := args[0]
	if len(args) > 1 {
		n, err = strconv.Atoi(args[1])
		if err != nil {
			return err
		}
	}

	totalTime := 0
	for i := 0; i < n; i++ {
		startTime := time.Now().Nanosecond()
		results, err := p.Proxy.Search(filename)
		if err != nil {
			return err
		}
		if results == nil {
			return errors.New("file not found " + filename)
		}
		endTime := time.Now().Nanosecond()

		totalTime += endTime - startTime
	}
	avgTime := int(totalTime / n)

	log.Info("Test search '%v' %v times: avg=%.2fms, total=%.2fms", filename, n, float64(avgTime)/float64(time.Millisecond), float64(totalTime)/float64(time.Millisecond))
	return nil
}
