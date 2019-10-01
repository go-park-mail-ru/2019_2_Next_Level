package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	Whitelist     map[string]bool
	StaticDir     string
	Port          string
	FileForFolder string
	OpenDir       string
	PrivateDir    string
	SelfURL       string
	DefaultAvatar string
}

func (c *Config) Inflate() error {
	filename := "config.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		return errors.New("Cannot open config.json")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &c)
	return nil
}
