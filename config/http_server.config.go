package config

import (
	"encoding/json"
)

type HTTPConfig struct {
	Config
	Whitelist     map[string]bool
	StaticDir     string
	Port          string
	FileForFolder string
	OpenDir       string
	PrivateDir    string
	SelfURL       string
	DefaultAvatar string
}

func (c *HTTPConfig) Inflate() error {
	filename := "http_server.config.json"
	var byteValue []byte
	err := c.readFile(filename, &byteValue)
	if err != nil {
		// fmt.Printf("Error during opening %s: %s", filename, err)
		return err
	}
	json.Unmarshal([]byte(byteValue), &c)
	return nil
}

// var Configuration HTTPConfig
