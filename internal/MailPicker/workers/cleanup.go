package workers

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"context"
	"fmt"
	gomail "github.com/veqryn/go-email/email"
//gomail "net/mail"
	"strings"
	"sync"
	"time"
)

type MailCleanup struct {
	errorChan chan error
}

func NewMailCleanup(errorChan chan error) *MailCleanup {
	return &MailCleanup{errorChan: errorChan}
}

func (w *MailCleanup) Run(externwg *sync.WaitGroup, ctx context.Context, in chan interface{}, out chan model.Email) {
	defer externwg.Done()

	go func() {
		for mail := range in {
			fmt.Println("Picker")
			emailTemp, ok := mail.(post.Email)
			if !ok {
				w.errorChan <- fmt.Errorf("Cannot convert input to email")
				return
			}
			email := model.Email{}
			email.Body = emailTemp.Body
			email.From = emailTemp.From
			email.To = emailTemp.To
			email.Body += "\n\n" // preventing the EOF error of bloody the parser
			reader := strings.NewReader(email.Body)
			msg, err := gomail.ParseMessage(reader)
			if err != nil {
				w.errorChan <- err
				return
			}

			if err := w.HandleMail(msg); err != nil {
				w.errorChan <- err
				return
			}

			res := w.Repack(msg)
			res.From = emailTemp.From
			res.To = emailTemp.To
			out <- res
			time.Sleep(20*time.Millisecond) // for tests
		}
	}()
	<-ctx.Done()
}

func (w *MailCleanup) HandleMail(email *gomail.Message) error{
	return nil
}

func (w *MailCleanup) Repack(from *gomail.Message) (model.Email) {
	to := model.Email{}
	to.Header.From = from.Header.From()
	to.Header.Subject = from.Header.Subject()
	//to.Header.WhenReceived, _ = from.Header.Date()
	to.Header.To = from.Header.To()
	if len(from.Parts)>0 {
		to.Body = string(from.Parts[0].Body)
	}else{
		to.Body = string(from.Body)
	}
	return to
}