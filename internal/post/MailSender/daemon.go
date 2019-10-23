package mailsender

import (
	"2019_2_Next_Level/internal/logger"
	"2019_2_Next_Level/internal/post"
	"sync"
)

// MailSender : checks emails, prepares them for delivery,
// controls if they are sent and reacts when the email cannot be sent
type MailSender struct {
	queueChan post.ChanPair
	smtpChan  post.ChanPair
	log       logger.Log
}

// Init : gets channel packs
func (s *MailSender) Init(pre, next post.ChanPair) error {
	s.queueChan = pre
	s.smtpChan = next
	s.log.SetPrefix("MailSender")
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
		// log.Println(email.Body)
		// log.Debug(email.Body)
		// s.log.Println(email.Body)
		s.smtpChan.Out <- email
		i++

	}
}
