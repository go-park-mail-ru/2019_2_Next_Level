package config

import (
	"fmt"
	"log"
)

type Config struct {
	RemotePort         string `json:"RemotePort"`
	RemoteHost         string `json:"RemoteHost"`
	RemoteCheckTimeout int16  `json:"remoteCheckTimeout"`
	PickerWorkerCount int `json:"pickerWorkerCount"`
	CleanerWorkerCount int`json:"cleanerWorkerCount`
	SaverWorkerCount int`json:"saverWorkerCount`
	DB Database `json:"Database"`
}
var defaultConfig Config
func (c *Config) Init(args ...interface{}) {
	//defaultConfig = Config{"1000","localhost","100",1,1,
	//	Database{
	//		Port:"1000",
	//		Host:"localhost",
	//
	//	}}
	if c.RemotePort != "" && c.RemotePort[0] != ':' {
		c.RemotePort = ":" + c.RemotePort
	}

	if len(args) >= 2{
		user, ok := args[0].(string)
		if !ok {
			log.Println("Wrong param for user in argument 0")
			return
		}
		pass, ok := args[1].(string)
		if !ok {
			log.Println("Wrong param for pass in argument 1")
			return
		}
		c.DB.Init(user, pass)
	} else {
		fmt.Println("Wrong param count", len(args))
		return
	}
}

type Database struct {
	Port string
	Host string
	DBName string
	User string
	Password string
}

func (d *Database) Init(user, pass string) {
	d.User = user
	d.Password = pass
}

var Conf Config