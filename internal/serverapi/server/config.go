package server

import (
	"2019_2_Next_Level/pkg/config"
)

type HTTPConfig struct {
	config.Config
	Whitelist           map[string]bool
	StaticDir           string
	Port                string
	FileForFolder       string
	OpenDir             string
	PrivateDir          string
	SelfURL             string
	DefaultAvatar       string
	PostServiceHost     string
	PostServiceSendPort string
}

func (c *HTTPConfig) Init() {
	c.Port = ":" + c.Port
	c.PostServiceSendPort = ":" + c.PostServiceSendPort
}

var Conf HTTPConfig
