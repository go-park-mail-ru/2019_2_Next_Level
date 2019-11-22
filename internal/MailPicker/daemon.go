package mailpicker

import (
	"2019_2_Next_Level/internal/MailPicker/config"
	"2019_2_Next_Level/internal/MailPicker/log"
	"2019_2_Next_Level/internal/MailPicker/workers"
	"2019_2_Next_Level/internal/model"
	postinterface "2019_2_Next_Level/internal/postInterface"
	e "2019_2_Next_Level/pkg/Error"
	"context"
	"fmt"
	"sync"
)

type Secretary struct {
	// postRes               postinterface.QueueClient
	postRes               postinterface.IPostInterface
	repo                  Repository
	queueConnectionStatus bool
	quitChan              chan interface{}
	errorChan chan error
}

func NewSecretary(postRes postinterface.IPostInterface, repo Repository, quitChan chan interface{}) *Secretary {
	return &Secretary{postRes: postRes, repo: repo, quitChan: quitChan}
}

// Init : initializes the module
func (s *Secretary) Init() *Secretary {
	s.postRes.Init()
	s.queueConnectionStatus = true
	s.errorChan = make(chan error, 3)

	log.Log().L(fmt.Sprintf("Init MailPicker listening remote %s%s", config.Conf.RemoteHost, config.Conf.RemotePort))
	return s
}

func (s *Secretary) Run(externwg *sync.WaitGroup) {
	defer externwg.Done()
	wg := sync.WaitGroup{}

	ctx1, finish1 := context.WithCancel(context.Background())
	ctx2, finish2 := context.WithCancel(context.Background())
	ctx3, finish3 := context.WithCancel(context.Background())
	chan1 := make(chan interface{}, 10)
	picker := workers.NewMailPicker(s.errorChan, s.postRes, s.repo.UserExists)
	for i:=0; i<config.Conf.PickerWorkerCount; i++ {
		wg.Add(1)
		go picker.Run(&wg, ctx1,chan1)
	}

	cleaner := workers.NewMailCleanup(s.errorChan)
	chan2 := make(chan model.Email, 10)
	for i:=0; i<config.Conf.CleanerWorkerCount; i++ {
		wg.Add(1)
		go cleaner.Run(&wg, ctx2, chan1, chan2)
	}

	saver := workers.NewMailSaver(s.errorChan, s.repo.AddEmail)
	for i:=0; i<config.Conf.SaverWorkerCount; i++ {
		wg.Add(1)
		go saver.Run(&wg, ctx3, chan2)
	}
	waitWgChan := make(chan interface{})
	go func(){
		wg.Wait()
		waitWgChan<-struct{}{}
	}()
	for {
		select {
		case err := <-s.errorChan:
			_, ok := err.(e.Error)
			if !ok{
				log.Log().E("Daemon stopping due to error: ", err)
				finish1()
				finish2()
				finish3()
				return
			}
			continue
			break;
		case <-waitWgChan:
			return
		}
			return
	}
	//wg.Wait()
}

