package client

import (
	"fmt"
	"net/rpc"
)

type Pool struct {
	clients map[string]*rpc.Client
}

func NewPool() *Pool {
	p := Pool{}
	p.clients = make(map[string]*rpc.Client)
	return &p
}

func (p *Pool) Dial(address string, port int) (cli *rpc.Client, err error) {
	key := fmt.Sprintf("%v:%v", address, port)

	cli, ok := p.clients[key]

	if ok {
		ok = p.test(cli)
	}

	if !ok {
		cli, err = rpc.Dial("tcp", key) // Connect to the other peer
		if err != nil {
			return nil, err
		}
		p.clients[key] = cli
	}

	return cli, nil
}

func (p *Pool) test(cli *rpc.Client) bool {
	var hello string = "hello"
	var reply string
	err := cli.Call("Handler.Echo", hello, &reply)
	if err != nil {
		return false
	}

	if reply != hello {
		return false
	}

	return true
}
