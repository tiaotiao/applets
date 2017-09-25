package main

import "./server"
import "./client"
import "../centralserver/proxy"
import "../common/commander"

type Peer struct {
	Server *server.Server
	Client *client.Client
	Proxy  *proxy.Proxy
	Cmd    *commander.Commander
}

func NewPeer(peerId string, port int, folderDir string, centralServerAddr string) *Peer {
	p := &Peer{}
	p.Proxy = proxy.NewProxy(centralServerAddr)
	p.Server = server.NewServer(port, peerId, p.Proxy)
	p.Client = client.NewClient(p.Proxy, folderDir)

	p.Cmd = commander.NewCommander()
	p.Cmd.Register("search", "[filename] ", p.CmdSearch)
	return p
}

func (p *Peer) CmdSearch(args ...string) error {
	for _, filename := range args {
		results, err := p.Proxy.Search(filename)
		if err != nil {
			return err
		}
		if results == nil {
			println("file not found", filename)
		} else {
			println(filename, results.FileInfo)
		}
	}
	return nil
}

func (p *Peer) Run() {
	go p.Server.Run()

	p.Cmd.Run()

	p.Server.Stop()
}
