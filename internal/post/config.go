package post

type PostServiceConfig struct {
	OutcomingQueue  MessageQueueConfig
	IncomingQueue   MessageQueueConfig
	Smtp            SMTPConfig
	ChannelCapasity int
}

func (c *PostServiceConfig) Init() {

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
