package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

func readParams(defaultFile string) (configfile string) {
	flag.StringVar(&configfile, "config", "", "Config file")
	flag.Parse()
	if configfile == "" {
		configfile = defaultFile
	}
	return
}

func LoadConfig(defaultFile string, v interface{}) error {
	file := readParams(defaultFile)

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}
