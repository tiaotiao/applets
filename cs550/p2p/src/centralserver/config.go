package main

import (
	"common/config"
	"errors"
	"flag"
)

type CentralServerConfig struct {
	Port  int  `json:"port"`
	Debug bool `json:"debug"`
}

func LoadConfig() (*CentralServerConfig, error) {
	cfg := CentralServerConfig{}

	file := readParams()
	if file == "" {
		return nil, errors.New("config file is not specified. \nUseage:\tcentralserver -config=CONFIG_FILE")
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
