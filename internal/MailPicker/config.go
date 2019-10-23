package mailpicker

type Config struct {
	RemotePort         string `json:"RemotePort"`
	RemoteHost         string `json:"RemoteHost"`
	RemoteCheckTimeout int16  `json:"remoteCheckTimeout"`
}

func (c *Config) Init() {

	if c.RemotePort != "" && c.RemotePort[0] != ':' {
		c.RemotePort = ":" + c.RemotePort
	}
}

var Conf Config
