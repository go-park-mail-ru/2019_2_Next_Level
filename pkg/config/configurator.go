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
	Init(...interface{})
}

func Inflate(filename string, dest ConfigInterface, args ...interface{}) error {
	jsonFile, err := os.Open(configPath + filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	output, err := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(output), &dest)
	(dest).Init(args...)
	return err
}

