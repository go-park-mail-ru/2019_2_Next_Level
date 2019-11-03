package config

import (
	"2019_2_Next_Level/internal/serverapi/server/config"
)

type MainConfig struct {
	HttpConfig config.HTTPConfig `json:"HttpServer"`
}

func (c *MainConfig) Init() {
	c.HttpConfig.Init()
	config.Conf = c.HttpConfig
}

var Conf MainConfig
