package sdk

import "net/rpc"
import "common"
import "errors"

type FileClient struct {
	address string
	rpc     *rpc.Client
}

func NewFileClient(addr string) *FileClient {
	c := &FileClient{}
	c.address = addr
	return c
}

func (c *FileClient) Connect() (err error) {
	c.rpc, err = rpc.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	return nil
}

//////////////////////////////////////////////////////////////////////

func (c *FileClient) Open(path string) (fileID string, size int, err error) {
	var result common.OpenResult
	err = c.rpc.Call("ClientHandler.Open", path, &result)
	if err != nil {
		return
	}
	fileID = result.FileID
	size = result.Size
	return
}

func (c *FileClient) Read(fileID string, offset int, n int) (data []byte, err error) {
	var param common.ReadParam
	param.FileID = fileID
	param.Offset = offset
	param.Count = n

	err = c.rpc.Call("ClientHandler.Read", param, &data)

	return
}

func (c *FileClient) Write(fileID string, data []byte) error {
	var param common.WriteParam
	param.FileID = fileID
	param.Data = data

	var n int
	err := c.rpc.Call("ClientHandler.Write", param, &n)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.New("n!=len(data)")
	}
	return nil
}

func (c *FileClient) AcquireRead(fileID string) (ok bool, err error) {
	err = c.rpc.Call("ClientHandler.AcquireRead", fileID, &ok)
	return
}

func (c *FileClient) AcquireWrite(fileID string) (ok bool, err error) {
	err = c.rpc.Call("ClientHandler.AcquireWrite", fileID, &ok)
	return
}

func (c *FileClient) Release(fileID string) (ok bool, err error) {
	err = c.rpc.Call("ClientHandler.Release", fileID, &ok)
	return
}

func (c *FileClient) Close(fileID string) (bool, error) {
	var ok bool
	err := c.rpc.Call("ClientHandler.Close", fileID, &ok)
	if err != nil {
		return false, err
	}
	return true, nil
}
