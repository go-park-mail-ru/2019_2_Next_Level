package smtpd

import (
	"2019_2_Next_Level/internal/logger"
	"2019_2_Next_Level/internal/post"

	"fmt"
	"sync"
)

// Server : class for SMTP server daemon
type Server struct {
	mailSenderChan    post.ChanPair
	incomingQueueChan post.ChanPair
	log               logger.Log
}

// Init : gets MailSender and IncomingQueue channels
func (s *Server) Init(pre, next post.ChanPair) error {
	s.mailSenderChan = pre
	s.incomingQueueChan = next
	fmt.Println("SMTPd started. Hello!")
	s.log.SetPrefix("SMTPd")
	return nil
}

// Run : start daemon's work
func (s *Server) Run(externwg *sync.WaitGroup) {
	defer externwg.Done()
	s.PrintAndForward()
}

// PrintAndForward : gets message from MailSender, prints it and resends to IncomingQueue
func (s *Server) PrintAndForward() {
	i := 0
	for pack := range s.mailSenderChan.Out {
		email := pack.(post.Email)
		s.log.Println(email.Body)
		s.incomingQueueChan.In <- email
		i++
	}
}
