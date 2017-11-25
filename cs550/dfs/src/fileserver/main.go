package main

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		println("LoadConfig error:", err.Error())
		return
	}

	s := NewFileServer(cfg.ExternalPort, cfg.InternalPort, cfg.LockServer, cfg.FileServers)

	err = s.Run()
	if err != nil {
		println("Run file server error:", err.Error())
		return
	}

	s.Stop()
}
