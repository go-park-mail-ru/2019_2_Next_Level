package incomingmailsecretary

type Config struct {
	RemotePort string `json:"RemotePort"`
	RemoteHost string `json:"RemoteHost"`
}

func (c *Config) Init() {

	if c.RemotePort != "" && c.RemotePort[0] != ':' {
		c.RemotePort = ":" + c.RemotePort
	}
}

var Conf Config
