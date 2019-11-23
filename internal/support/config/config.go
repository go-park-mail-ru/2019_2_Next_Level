package config

import (
	"fmt"
	"log"
)

type MainConfig struct {
	HttpConfig HttpConfig `json:"HttpServer"`
	DB Database `json:"Database"`
}

func (c *MainConfig) Init(args ...interface{}) {
	//c.HttpConfig.Init()
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

type HttpConfig struct {

}

var Conf MainConfig
