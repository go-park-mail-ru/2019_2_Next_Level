package serverapi

import (
	"2019_2_Next_Level/internal/serverapi/server"
)

type MainConfig struct {
	HttpConfig server.HTTPConfig `json:"HttpServer"`
}

func (c *MainConfig) Init() {
	c.HttpConfig.Init()
	server.Conf = c.HttpConfig
}

var Conf MainConfig
