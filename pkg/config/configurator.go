package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	configPath = ""
)

type ConfigInterface interface {
	Init(...interface{})
}

func Inflate(filename string, dest ConfigInterface, args ...interface{}) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location := filepath.Dir(ex)
	fmt.Println(location)
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

