package client

import (
	"errors"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"

	"../../centralserver/proxy"
	"../../common"
)

type Client struct {
	centralServer *proxy.Proxy
	folder        string
}

func NewClient(p *proxy.Proxy, folderPath string) *Client {
	c := &Client{}
	c.centralServer = p
	return c
}

func (c *Client) RequestFile(fileName string) error {
	results, err := c.centralServer.Search(fileName)
	if err != nil {
		return err
	}

	if len(results.Peers) <= 0 {
		return errors.New("peers list is empty")
	}

	// TODO choose a peer
	p := results.Peers[0]

	err = c.obtain(p, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) obtain(p common.PeerInfo, fileName string) error {
	rpcClient, err := rpc.Dial("tcp", p.Address)
	if err != nil {
		return err
	}
	var content []byte
	err = rpcClient.Call("Obtain", fileName, &content)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(c.folder, fileName), content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
