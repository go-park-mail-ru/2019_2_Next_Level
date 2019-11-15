package config

type HTTPConfig struct {
	//config.Config
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
	HostName string `json:"hostname"`
}

func (c *HTTPConfig) Init(args ...interface{}) {
	c.Port = ":" + c.Port
	c.PostServiceSendPort = ":" + c.PostServiceSendPort
}

var Conf HTTPConfig
