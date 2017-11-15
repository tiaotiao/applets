package main

import (
	"common/log"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

///////////////////////////////////////////////
// Commands

func (p *Peer) cmdAddFile(args ...string) error {
	for _, filepath := range args {
		err := p.Server.FileMgr.AddFile(filepath)
		if err != nil {
			return err
		}
		log.Info("Add file ok: %s", filepath)
	}
	return nil
}

func (p *Peer) cmdSearch(args ...string) error {
	for _, filename := range args {
		results, err := p.Proxy.Search(filename)
		if err != nil {
			return err
		}
		if !results.Exist {
			log.Info("Search '%s' not found", filename)
		} else {
			log.Info("Search '%s' size=%v, peers=%v", filename, results.Size, results.Peers)
		}
	}
	return nil
}

func (p *Peer) cmdListAll(args ...string) error {
	results, err := p.Proxy.ListAll()
	if err != nil {
		return err
	}
	log.Info("List %d files", len(results))
	var i = 1
	for _, r := range results {
		log.Info("%-5d %-15vsize=%-8v peers=%v", i, r.Name, r.Size, r.StringPeers())
		i++
	}
	return nil
}

func (p *Peer) cmdObtain(args ...string) error {
	for _, filename := range args {
		err := p.Client.Obtain(filename) // print log inside
		if err != nil {
			return err
		}
		log.Info("[Cmd] Obtain file %v ok.", filename)
	}
	return nil
}

func (p *Peer) cmdTestProformance(args ...string) (err error) {
	n := 10000
	if len(args) < 0 {
		return nil
	}
	if len(args) > 0 {
		n, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
	}

	now := time.Now()
	delaySec := 5 - now.Second()%5
	if delaySec < 2 {
		delaySec += 5
	}
	log.Info("Delay %v seconds...", delaySec)
	<-time.After(time.Duration(delaySec) * time.Second) // Delay few seconds to wait for other peers

	log.Info("Start testing %v", n)

	filelist, err := p.Proxy.ListAll()
	if err != nil {
		return err
	}

	randFilename := func() string {
		idx := rand.Intn(len(filelist))
		filename := ""
		i := 0
		for name := range filelist { // choose a random file to search
			if i == idx {
				filename = name
			}
			i++
		}
		return filename
	}

	var logTime int64 = 0
	var totalTime int64 = 0
	for i := 0; i < n; i++ {
		startTime := time.Now().UnixNano()

		prob := rand.Intn(100)

		if prob < 10 { // registry 10%
			err = p.Server.FileMgr.AddFolder(p.Dir)
			if err != nil {
				//return err
				log.Error("Registy error: %v", err)
			}

		} else if prob < 90 { // search 80%
			filename := randFilename()
			_, err := p.Proxy.Search(filename) // Search file
			if err != nil {
				//return err
				log.Error("Search error: %v", err)
			}

		} else { // obtain 10%
			filename := randFilename()
			err := p.Client.Obtain(filename)
			if err != nil {
				//return err
				log.Error("Obtain error: %v", err)
			}
		}

		endTime := time.Now().UnixNano()
		elapsedTime := endTime - startTime

		logTime += elapsedTime
		totalTime += elapsedTime

		if time.Duration(logTime) > time.Second*5 {
			log.Info("Testing %v/%v ...", i, n)
			logTime = 0
		}
	}
	avgTime := totalTime / int64(n)

	log.Info("Test ends %v: avg=%.2fms, total=%.2fms", n, float64(avgTime)/float64(time.Millisecond), float64(totalTime)/float64(time.Millisecond))

	return nil
}
