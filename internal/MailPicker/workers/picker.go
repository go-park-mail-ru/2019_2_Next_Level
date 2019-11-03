package workers

import (
	"2019_2_Next_Level/internal/MailPicker/config"
	postinterface "2019_2_Next_Level/internal/postInterface"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type MailPicker struct {
	errorChan chan error
	input postinterface.IPostInterface
	checkUserExist func(string) bool
	inputStatus bool
}

func NewMailPicker(errorChan chan error, input postinterface.IPostInterface, checkUserExist func(string) bool) *MailPicker {
	return &MailPicker{errorChan: errorChan, input: input, checkUserExist: checkUserExist}
}


func (w *MailPicker) Run(externwg *sync.WaitGroup, ctx context.Context, out chan interface{}) {
	w.inputStatus = true
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			email, err := w.input.Get()
			if err != nil {
				if w.inputStatus == true {
					w.inputStatus = false
					fmt.Println(err)
				}
				time.Sleep(time.Duration(config.Conf.RemoteCheckTimeout) * time.Millisecond)
				continue
			}
			if !w.inputStatus {
				fmt.Println("Connection emerged!")
				w.inputStatus = !w.inputStatus
			}
			fmt.Println("Messsage got")
			if w.checkUserExist(strings.Split(email.To, "@")[0]) {
				out <- email
				time.Sleep(20*time.Millisecond) // for tests
			}
		}
	}()
	<-ctx.Done()
}
