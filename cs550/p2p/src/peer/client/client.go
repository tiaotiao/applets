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

	err = c.verifyFile(&results.FileInfo)
	if err != nil {
		return err
	}

	// TODO register to central server

	log.Info("Download file ok. '%v' size=%v, content=[%v], md5=%v", results.Name, results.Size, snippet, results.Md5)

	return nil
}

// Download file from the other peer
func (c *Client) obtain(p common.PeerInfo, fileName string) (string, error) {
	rpcClient, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", p.Address, p.Port)) // Connect to the other peer
	if err != nil {
		return "", err
	}

	// WAENING: Due to the limitation of the RPC library, this RPC call will load the whole file in memory before
	// 		write to file. It will be a big issue if the file is very large. This function need to be changed later.
	var content []byte                                         // File content will be loaded in here
	err = rpcClient.Call("Handler.Obtain", fileName, &content) // RPC call
	if err != nil {
		return "", err
	}

	snippet := "" // snippet of the content
	if len(content) <= 64 {
		snippet = string(content)
	} else {
		snippet = string(content[:30]) + "...." + string(content[len(content)-30:])
	}
	snippet = strings.TrimSpace(snippet)

	path := filepath.Join(c.folder, fileName)
	err = ioutil.WriteFile(path, content, os.ModePerm) // Write content to file
	if err != nil {
		log.Error("Write file error %v, %v, %v", err, path, len(content))
		return snippet, err
	}

	return snippet, nil
}

// Verify downloaded file by size and md5
func (c *Client) verifyFile(f *common.FileInfo) error {
	path := filepath.Join(c.folder, f.Name)

	info, err := common.GetFileInfo(path) // Read file info
	if err != nil {
		return err
	}

	// Compare the size and md5
	if f.Size != info.Size || f.Md5 != info.Md5 {
		return fmt.Errorf("File not match remote=(%v,%v), local=(%v,%v) %v", f.Size, f.Md5, info.Size, info.Md5, info.Path)
	}

	return nil
}
