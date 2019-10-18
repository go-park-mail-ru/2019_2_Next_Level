package config

import (
	"encoding/json"
)

type OutpqConfig struct {
	Config
	Host string
	Port string
}

func (c *OutpqConfig) Inflate() error {
	filename := "outcoming_queue.config.json"
	var byteValue []byte
	err := c.readFile(filename, &byteValue)
	if err != nil {
		// fmt.Printf("Error during opening %s: %s", filename, err)
		return err
	}
	json.Unmarshal([]byte(byteValue), &c)
	return nil
}

// var Configuration Config
