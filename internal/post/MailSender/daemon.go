package mailsender

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/log"
	"sync"
)

// MailSender : checks emails, prepares them for delivery,
// controls if they are sent and reacts when the email cannot be sent
type MailSender struct {
	queueChan post.ChanPair
	smtpChan  post.ChanPair
}

// Init : gets channel packs
func (s *MailSender) Init(pre, next post.ChanPair, _ ...interface{}) error {
	s.queueChan = pre
	s.smtpChan = next
	log.Log().L("Init mailsender")
	return nil

}

// Run : starts the daemon's work
func (s *MailSender) Run(externWg *sync.WaitGroup) {
	defer externWg.Done()
	s.ProcessEmail()

}

// ProcessEmail : handles messages from the queue fro sent
func (s *MailSender) ProcessEmail() {
	i := 0
	for pack := range s.queueChan.Out {
		email := pack.(post.Email)
		s.smtpChan.Out <- email
		i++
	}
}
