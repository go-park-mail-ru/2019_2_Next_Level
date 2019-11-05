package workers

import (
	"2019_2_Next_Level/internal/model"
	e "2019_2_Next_Level/pkg/Error"
	"context"
	"sync"
	"time"
)

type MailSaver struct {
	errorChan chan error
	saveMailFunc func(email *model.Email) error
}

func NewMailSaver(errorChan chan error, saveMailFunc func(email *model.Email) error) *MailSaver {
	return &MailSaver{errorChan: errorChan, saveMailFunc: saveMailFunc}
}

func (w *MailSaver) Run(externwg *sync.WaitGroup, ctx context.Context, in chan model.Email) {
	defer externwg.Done()
	go func() {
			for email := range in {
			w.ProcessEmail(&email)
			time.Sleep(20*time.Millisecond) // for tests
		}
	}()
	<-ctx.Done()
}

func (w *MailSaver) ProcessEmail(email *model.Email) error {
	err := w.saveMailFunc(email)
	if err != nil {
		w.errorChan <- e.Error{}.SetCode(e.ProcessError).SetError(err).SetPlace("MailSaver");
	}
	return nil
}

