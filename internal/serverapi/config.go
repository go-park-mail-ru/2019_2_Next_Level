package serverapi

import (
	incommail "2019_2_Next_Level/internal/serverapi/IncomingMailSecretary"
	"2019_2_Next_Level/internal/serverapi/server"
)

type MainConfig struct {
	HttpConfig        server.HTTPConfig `json:"HttpServer"`
	IncomingSecretary incommail.Config  `json:"incomingSecretary"`
}

func (c *MainConfig) Init() {
	c.HttpConfig.Init()
	server.Conf = c.HttpConfig

	c.IncomingSecretary.Init()
	incommail.Conf = c.IncomingSecretary
}

// type HTTPConfig struct {
// 	config.Config
// 	Whitelist           map[string]bool
// 	StaticDir           string
// 	Port                string
// 	FileForFolder       string
// 	OpenDir             string
// 	PrivateDir          string
// 	SelfURL             string
// 	DefaultAvatar       string
// 	PostServiceHost     string
// 	PostServiceSendPort string
// }

// func (c *HTTPConfig) Init() {
// 	c.Port = ":" + c.Port
// 	c.PostServiceSendPort = ":" + c.PostServiceSendPort
// }

// var Conf HTTPConfig
var Conf MainConfig
