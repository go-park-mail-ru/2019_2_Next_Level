package workers

import (
	"2019_2_Next_Level/internal/MailPicker/log"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"context"
	"fmt"
	gomail "github.com/veqryn/go-email/email"
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
			email.Header.Subject = emailTemp.Subject
			var res model.Email

			if email.From != "mailder-daemon@nl-mail.ru" {
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
				res = w.Repack(msg)
				log.Log().L(res)
			} else {
				res.Body = email.Body
			}
			res.From = emailTemp.From
			res.To = strings.Split(emailTemp.To, "@")[0]
			out <- res
			time.Sleep(200*time.Millisecond) // for tests
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
	to.Header.WhenReceived, _ = from.Header.Date()
	to.Header.To = from.Header.To()
	//for i, label := range to.Header.To {
	//	if len(label)<3 {
	//		continue
	//	}
	//	to.Header.To[i] = label[strings.Index(label, "<")+1 : strings.Index(label, ">")]
	//}
	to.Body = string(w.SelectBody(from))
	return to
}

func (w *MailCleanup) SelectBody(mail *gomail.Message) []byte {
	if len(mail.Parts) == 0{
		return mail.Body
	}
	body := make([]byte, 0)
	for _, part := range mail.MessagesAll() {
		mediaType, _, _ := part.Header.ContentType()
		switch mediaType {
		case "text/plain":
			if len(body) == 0 {
				body = part.Body
			}
			break
		case "text/html":
			body = part.Body
			break
		default:
			break
		}
	}
	return body
}