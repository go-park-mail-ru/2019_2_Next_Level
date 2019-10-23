package mailpicker

import (
	"2019_2_Next_Level/internal/MailPicker/repository"
	"2019_2_Next_Level/internal/model"
	postinterface "2019_2_Next_Level/internal/postInterface"
	"fmt"
	"sync"
	"time"
)

type Secretary struct {
	// postRes               postinterface.QueueClient
	postRes               postinterface.IPostInterface
	repo                  Repository
	queueConnectionStatus bool
	quitChan              chan interface{}
}

func NewInstanceDefault(dbConnection *model.Connection) Secretary {
	s := Secretary{}
	repo := repository.NewRepository(dbConnection)
	qInterface := postinterface.QueueClient{
		RemoteHost: Conf.RemoteHost,
		RemotePort: Conf.RemotePort,
	}
	s.Init(
		&repo,
		&qInterface,
	)
	return s
}
func (s *Secretary) DefaultInit(dbConnection *model.Connection) {
	repo := repository.NewRepository(dbConnection)
	qInterface := postinterface.QueueClient{
		RemoteHost: Conf.RemoteHost,
		RemotePort: Conf.RemotePort,
	}
	s.Init(
		&repo,
		&qInterface,
	)
}

func (s *Secretary) Init(repo Repository, qInterface postinterface.IPostInterface) error {
	// s.postRes = postinterface.QueueClient{RemoteHost: Conf.RemoteHost, RemotePort: Conf.RemotePort}
	s.postRes = qInterface
	s.postRes.Init()
	s.queueConnectionStatus = true

	s.repo = repo
	fmt.Println("Init MailPicker")
	return nil
}
func (s *Secretary) Run(externwg *sync.WaitGroup) {
	defer externwg.Done()
	s.MailPicker()
}

func (s *Secretary) MailPicker() {
	for {
		email, err := s.postRes.Get()
		if err != nil {
			if s.queueConnectionStatus == true {
				s.queueConnectionStatus = false
				fmt.Println(err)
			}
			time.Sleep(time.Duration(Conf.RemoteCheckTimeout) * time.Millisecond)
			continue
		}
		if !s.queueConnectionStatus {
			fmt.Println("Connection emerged!")
			s.queueConnectionStatus = !s.queueConnectionStatus
		}

		if !s.repo.UserExists(email.To) {
			fmt.Println("User not exists")
			continue
		}
		s.repo.AddEmail(model.Email(email))
	}
}
