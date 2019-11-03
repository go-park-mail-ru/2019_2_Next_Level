package smtpd

import (
	"2019_2_Next_Level/internal/logger"
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/smtpd/worker"
	"time"

	"fmt"
	"sync"

	"github.com/emersion/go-smtp"
)

type IncomingSmtpInerface interface {
	ListenAndServe() error
	Init(string, string) error
}

// Server : class for SMTP server daemon
type Server struct {
	mailSenderChan    post.ChanPair
	incomingQueueChan post.ChanPair
	log               logger.Log
	worker            worker.Worker
	smtpServer        IncomingSmtpInerface
	resultChannel     chan worker.EmailNil
	quitChan          chan interface{}
}

// Init : gets MailSender and IncomingQueue channels
func (s *Server) Init(pre, next post.ChanPair, args ...interface{}) error {
	// process the only iteration
	if len(args) == 1 {
		if res, ok := args[0].(IncomingSmtpInerface); ok {
			s.smtpServer = res
		} else {
			return fmt.Errorf("Wrong smtpServer got")
		}
	} else {
		s.smtpServer = s.NewDefaultSMTP("1025", "0.0.0.0")
	}
	s.resultChannel = make(chan worker.EmailNil, 100)
	s.quitChan = make(chan interface{}, 5)
	s.worker = worker.Worker{OutChannel: s.resultChannel}

	s.mailSenderChan = pre
	s.incomingQueueChan = next

	fmt.Println("SMTPd started. Hello!")
	s.log.SetPrefix("SMTPd")

	return nil
}

func (s *Server) NewDefaultSMTP(port, host string) IncomingSmtpInerface {
	ss := &SMTPIncoming{smtp.NewServer(&s.worker)}
	ss.Init(port, host)
	return ss
}

// Run : start daemon's work
func (s *Server) Run(externwg *sync.WaitGroup) {
	defer externwg.Done()
	go s.RunSmtpServer()
	go s.GetIncomingMessages()
	// go s.GenAndSendMailTest()
	// go s.PrintAndForward()
	select {
	case <-s.quitChan:
		return
	}
}

func (s *Server) GenAndSendMailTest() {
	for {
		email := post.Email{"ivan", "ian", `Received: from mxback18o.mail.yandex.net (mxback18o.mail.yandex.net [IPv6:2a02:6b8:0:1a2d::69])
			by forward102o.mail.yandex.net (Yandex) with ESMTP id 1237E6680F6C
			From: Andrey K. <andreykochnov@yandex.ru>
			To: aaa <aaa@nlmail.ddns.net>
			Subject: Test
			MIME-Version: 1.0
			Date: Sun, 03 Nov 2019 20:54:17 +0300
			Content-Transfer-Encoding: 7bit
			Content-Type: text/html
			
			<div>Hellp</div>`}
		s.incomingQueueChan.In <- email
		time.Sleep(2000 * time.Millisecond)
	}
}

func (s *Server) RunSmtpServer() {
	// log.Println("Starting server at", s.smtpServer.Domain+s.smtpServer.Addr)
	if err := s.smtpServer.ListenAndServe(); err != nil {
		// s.log.Println("Error: cannot start incoming smtpServer")
		s.quitChan <- struct{}{}
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
		// s.log.Println(email.Body)
		s.incomingQueueChan.In <- email
		i++
	}
}

type SMTPIncoming struct {
	*smtp.Server
}

func (s *SMTPIncoming) Init(port, host string) error {
	s.Addr = ":" + port
	s.Domain = host
	s.ReadTimeout = 60 * time.Second
	s.WriteTimeout = 60 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	return nil
}
