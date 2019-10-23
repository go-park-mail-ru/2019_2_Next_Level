package smtpd

import (
	"2019_2_Next_Level/internal/logger"
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/smtpd/worker"
	"log"
	"time"

	"fmt"
	"sync"

	"github.com/emersion/go-smtp"
)

// Server : class for SMTP server daemon
type Server struct {
	mailSenderChan    post.ChanPair
	incomingQueueChan post.ChanPair
	log               logger.Log
	worker            worker.Worker
	smtpServer        *smtp.Server
	resultChannel     chan worker.EmailNil
}

// Init : gets MailSender and IncomingQueue channels
func (s *Server) Init(pre, next post.ChanPair) error {
	s.resultChannel = make(chan worker.EmailNil, 100)
	s.worker = worker.Worker{OutChannel: s.resultChannel}
	s.mailSenderChan = pre
	s.incomingQueueChan = next

	s.smtpServer = smtp.NewServer(&s.worker)

	s.smtpServer.Addr = ":" + "1025"
	s.smtpServer.Domain = "0.0.0.0"
	s.smtpServer.ReadTimeout = 60 * time.Second
	s.smtpServer.WriteTimeout = 60 * time.Second
	s.smtpServer.MaxMessageBytes = 1024 * 1024
	s.smtpServer.MaxRecipients = 50
	s.smtpServer.AllowInsecureAuth = true

	fmt.Println("SMTPd started. Hello!")
	s.log.SetPrefix("SMTPd")

	// incomingworker.outChan = next
	return nil
}

// Run : start daemon's work
func (s *Server) Run(externwg *sync.WaitGroup) {
	defer externwg.Done()
	go s.RunSmtpServer()
	go s.GetIncomingMessages()
	// go s.GenAndSendMailTest()
	s.PrintAndForward()
}

func (s *Server) GenAndSendMailTest() {
	for {
		email := post.Email{"ivan", "ian", "body"}
		s.incomingQueueChan.In <- email
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *Server) RunSmtpServer() {
	log.Println("Starting server at", s.smtpServer.Domain+s.smtpServer.Addr)
	if err := s.smtpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) GetIncomingMessages() {
	for data := range s.resultChannel {
		if data.Error != nil {
			fmt.Println("Wrong email getting: ", data.Error)
			continue
		}
		fmt.Println("Got message")
		fmt.Println(data.Email.Stringify())
		s.incomingQueueChan.In <- data.Email

	}
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
