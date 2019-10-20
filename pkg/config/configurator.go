package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	configPath = ""
)

type ConfigInterface interface {
	Init()
}

type Config struct {
}

func (c *Config) Inflate(filename string, dest ConfigInterface) error {
	jsonFile, err := os.Open(configPath + filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	output, err := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(output), &dest)
	(dest).Init()
	return err
}

var Configurator Config
