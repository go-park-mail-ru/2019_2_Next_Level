package post

type PostServiceConfig struct {
	OutcomingQueue  MessageQueueConfig
	IncomingQueue   MessageQueueConfig
	Smtp            SMTPConfig
	ChannelCapasity int
	Login string
	Password string
	Host string
	Port string
}

func (c *PostServiceConfig) Init(args ...interface{}) {

}

type MessageQueueConfig struct {
	Port            string
	ChannelCapasity int
}

type SMTPConfig struct {
	Port            string
	ChannelCapasity int
}

var Conf PostServiceConfig
