package main

import (
	"common/config"
	"errors"
	"flag"
)

type PeerConfig struct {
	Servers []string `json:"servers"`
	Dir     string   `json:"dir"`
	PeerID  string   `json:"id"`
	Port    int      `json:"port"`
	Debug   bool     `json:"debug"`
}

func LoadConfig() (*PeerConfig, error) {
	cfg := PeerConfig{}

	file := readParams()
	if file == "" {
		return nil, errors.New("config file is not specified. \nUseage:\tpeer -config=CONFIG_FILE")
	}

	err := config.LoadConfig(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func readParams() (configfile string) {
	flag.StringVar(&configfile, "config", "", "Config file.")

	flag.Parse()
	return
}
