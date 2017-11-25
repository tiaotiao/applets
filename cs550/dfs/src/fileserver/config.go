package main

import "common/config"
import "errors"

type Config struct {
	ExternalPort int      `json:"external_port"`
	InternalPort int      `json:"internal_port"`
	LockServer   string   `json:"lock_server"`
	FileServers  []string `json:"file_servers"`
}

func LoadConfig() (*Config, error) {
	c := Config{}

	err := config.LoadConfig("", &c)
	if err != nil {
		return nil, err
	}

	if c.ExternalPort == 0 {
		return nil, errors.New("missing external port")
	}
	if c.InternalPort == 0 {
		return nil, errors.New("missing internal port")
	}
	if c.LockServer == "" {
		return nil, errors.New("missing lock server address")
	}
	if len(c.FileServers) == 0 {
		return nil, errors.New("missing file server addresses")
	}

	return &c, nil
}
