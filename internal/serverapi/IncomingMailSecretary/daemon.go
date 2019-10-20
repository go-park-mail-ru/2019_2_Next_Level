package incomingmailsecretary

import (
	postinterface "2019_2_Next_Level/internal/serverapi/postInterface"
	"fmt"
	"sync"
	"time"
)

type Secretary struct {
	postRes postinterface.QueueClient
}

func (s *Secretary) Init() {
	s.postRes = postinterface.QueueClient{RemoteHost: Conf.RemoteHost, RemotePort: Conf.RemotePort}
	s.postRes.Init()
}
func (s *Secretary) Run(externwg *sync.WaitGroup) {
	defer externwg.Done()
	s.MailPicker()
}

func (s *Secretary) MailPicker() {
	for {
		email, err := s.postRes.Get()
		if err != nil {
			fmt.Println(err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		fmt.Println(email.Stringify())
	}
}
