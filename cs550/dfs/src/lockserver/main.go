package main

import "common"

func main() {
	s := NewLockServer(common.LOCKSERVER_PORT)
	s.Run()
	s.Stop()
}
