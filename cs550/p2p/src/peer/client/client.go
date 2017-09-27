package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"

	"centralserver/proxy"
	"common"
	"common/log"
)

type Client struct {
	centralServer *proxy.Proxy
	folder        string
}

func NewClient(p *proxy.Proxy, folderPath string) *Client {
	c := &Client{}
	c.centralServer = p
	c.folder = folderPath
	return c
}

func (c *Client) Obtain(fileName string) error {
	results, err := c.centralServer.Search(fileName)
	if err != nil {
		return err
	}

	if len(results.Peers) <= 0 {
		return errors.New("peers list is empty")
	}

	// TODO choose a peer
	p := results.Peers[0]

	snippet, err := c.obtain(p, fileName)
	if err != nil {
		return err
	}

	err = c.checkFile(&results.FileInfo)
	if err != nil {
		return err
	}

	log.Info("[Peer] Obtain file ok. %v size=%v md5=%v content=[%v]", results.Name, results.Size, results.Md5, snippet)

	return nil
}

func (c *Client) obtain(p common.PeerInfo, fileName string) (string, error) {
	rpcClient, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", p.Address, p.Port))
	if err != nil {
		return "", err
	}
	var content []byte
	err = rpcClient.Call("Handler.Obtain", fileName, &content)
	if err != nil {
		return "", err
	}

	snippet := ""
	if len(content) <= 64 {
		snippet = string(content)
	} else {
		snippet = string(content[:30]) + "...." + string(content[len(content)-30:])
	}

	snippet = strings.TrimSpace(snippet)

	path := filepath.Join(c.folder, fileName)
	err = ioutil.WriteFile(path, content, os.ModePerm)
	if err != nil {
		log.Debug("Write file error %v, %v, %v", err, path, len(content))
		return snippet, err
	}

	return snippet, nil
}

func (c *Client) checkFile(f *common.FileInfo) error {
	path := filepath.Join(c.folder, f.Name)

	info, err := common.GetFileInfo(path)
	if err != nil {
		return err
	}

	if f.Size != info.Size || f.Md5 != info.Md5 {
		return fmt.Errorf("file not match remote=(%v,%v), local=%v,%v %v", f.Size, f.Md5, info.Size, info.Md5, info.Path)
	}

	return nil
}
