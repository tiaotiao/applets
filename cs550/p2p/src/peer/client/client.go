package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"centralserver/proxy"
	"common"
)

type Client struct {
	centralServer *proxy.Proxy
	folder        string
	peerId        string
	pool          *Pool
}

func NewClient(peerId string, p *proxy.Proxy, folderPath string) *Client {
	c := &Client{}
	c.centralServer = p
	c.folder = folderPath
	c.peerId = peerId
	c.pool = NewPool()
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

	// choose a peer randomly
	idx := rand.Intn(len(results.Peers))
	p := results.Peers[idx]

	if p.PeerId == c.peerId {
		return nil
	}

	_, err = c.obtain(p, fileName)
	if err != nil {
		return err
	}

	// err = c.verifyFile(&results.FileInfo)
	// if err != nil {
	// 	return err
	// }

	// register this file to central server
	// report that the file is avaliable from this peer
	_, err = c.centralServer.Registry(&results.FileInfo)
	if err != nil {
		return err
	}

	//log.Info("Download file ok. '%v' size=%v, content=[%v], md5=%v", results.Name, results.Size, snippet, results.Md5)

	return nil
}

// Download file from the other peer
func (c *Client) obtain(p common.PeerInfo, fileName string) (string, error) {
	rpcClient, err := c.pool.Dial(p.Address, p.Port) // Connect to the other peer
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

	err = c.writeFile(fileName, content) // Write content to file
	if err != nil {
		//log.Error("Write file error %v, %v", err, fileName)
		return snippet, err
	}

	return snippet, nil
}

func (c *Client) writeFile(fileName string, content []byte) error {
	randFile := filepath.Join(c.folder, common.RandString(16))

	err := ioutil.WriteFile(randFile, content, os.ModePerm)
	if err != nil {
		return err
	}

	path := filepath.Join(c.folder, fileName)

	err = os.Rename(randFile, path)
	if err != nil {
		e := os.Remove(randFile)
		if e != nil {
			return e
		}
		return err
	}

	return nil
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
