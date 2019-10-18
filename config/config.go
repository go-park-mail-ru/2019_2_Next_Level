package config

import (
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
}

func (c *Config) readFile(filename string, output *[]byte) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return errors.New("Cannot open config.json")
	}
	defer jsonFile.Close()
	*output, err = ioutil.ReadAll(jsonFile)
	// json.Unmarshal([]byte(byteValue), &c)
	return err
}

var Configuration Config
