package main

import (
	"common/commander"
	"common/log"
	"fileserver/sdk"
)

type Client struct {
	fileserver *sdk.FileClient
	cmd        *commander.Commander

	path   string
	fileID string
}

func NewClient(serverAddr string) *Client {
	c := &Client{}
	c.fileserver = sdk.NewFileClient(serverAddr)
	c.cmd = commander.NewCommander()

	c.cmd.Register("open", c.CmdOpen, "open [filepath]")
	c.cmd.Register("read", c.CmdRead, "read from file")
	c.cmd.Register("write", c.CmdWrite, "write [data] to file")
	c.cmd.Register("acquire", c.CmdAcquire, "acquire [read|write] lock")
	c.cmd.Register("release", c.CmdRelease, "release lock")
	c.cmd.Register("close", c.CmdClose, "close file")
	return c
}

func (c *Client) Run() error {
	err := c.fileserver.Connect()
	if err != nil {
		return err
	}

	c.cmd.Run() // block

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func (c *Client) CmdOpen(args ...string) error {
	if len(args) < 1 {
		return nil
	}
	path := args[0]

	fileID, size, err := c.fileserver.Open(path)
	if err != nil {
		return err
	}
	c.path = path
	c.fileID = fileID
	log.Info("Open file %v: fileID=%v, size=%v", path, fileID, size)
	return nil
}

func (c *Client) CmdRead(args ...string) error {
	if c.fileID == "" {
		log.Error("File not open")
		return nil
	}

	data, err := c.fileserver.Read(c.fileID, 0, -1)
	if err != nil {
		return err
	}

	log.Info("Read file: %v", string(data))
	return nil
}

func (c *Client) CmdWrite(args ...string) error {
	if len(args) < 1 {
		return nil
	}

	data := "\n" + args[0]

	if c.fileID == "" {
		log.Error("File not open")
		return nil
	}

	err := c.fileserver.Write(c.fileID, []byte(data))
	if err != nil {
		return err
	}

	log.Info("Write file OK")
	return nil
}

func (c *Client) CmdAcquire(args ...string) (err error) {
	if len(args) < 1 {
		return nil
	}
	if c.fileID == "" {
		log.Error("File not open")
		return nil
	}
	mode := args[0]

	var ok bool
	if mode == "read" {
		ok, err = c.fileserver.AcquireRead(c.fileID)
	} else if mode == "write" {
		ok, err = c.fileserver.AcquireWrite(c.fileID)
	} else {
		log.Error("%v invalid, should be 'read' or 'write'", mode)
	}

	if ok {
		log.Info("Acquire lock OK.")
	} else {
		log.Info("Acquire lock FAILED.")
	}

	return nil
}

func (c *Client) CmdRelease(args ...string) (err error) {
	if c.fileID == "" {
		log.Error("File not open")
		return nil
	}

	var ok bool
	ok, err = c.fileserver.Release(c.fileID)

	if ok {
		log.Info("Release lock OK.")
	} else {
		log.Info("Release lock FAILED.")
	}
	return nil
}

func (c *Client) CmdClose(args ...string) error {
	if c.fileID == "" {
		log.Error("File not open")
		return nil
	}

	ok, err := c.fileserver.Close(c.fileID)
	if err != nil {
		return err
	}

	log.Info("Close file %v", ok)
	c.fileID = ""
	c.path = ""
	return nil
}
