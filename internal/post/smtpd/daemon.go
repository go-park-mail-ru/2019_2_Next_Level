package smtpd

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/log"
	"2019_2_Next_Level/internal/post/smtpd/worker"
	"bytes"
	"fmt"
	gomail "gopkg.in/gomail.v2"
	"sync"
	"time"

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
	worker            worker.Worker
	smtpServer        IncomingSmtpInerface
	resultChannel     chan worker.EmailNil
	quitChan          chan interface{}
}

// Init : gets MailSender and IncomingQueue channels
func (s *Server) Init(pre, next post.ChanPair, args ...interface{}) error {
	s.resultChannel = make(chan worker.EmailNil, 100)
	// process the only iteration
	if len(args) == 1 {
		if res, ok := args[0].(IncomingSmtpInerface); ok {
			s.smtpServer = res
		} else {
			return fmt.Errorf("Wrong smtpServer got")
		}
	} else {
		s.smtpServer = s.NewDefaultSMTP(post.Conf.ListenPort, "0.0.0.0")
	}
	//s.resultChannel = make(chan worker.EmailNil, 100)
	s.quitChan = make(chan interface{}, 5)
	s.worker = worker.Worker{OutChannel: s.resultChannel}

	s.mailSenderChan = pre
	s.incomingQueueChan = next

	notification := "SMTPd started. Port to listen: %s"
	log.Log().I(fmt.Sprintf(notification, post.Conf.ListenPort))

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
	log.Log().L("Run SMTP...")
	go s.RunSmtpServer()
	go s.GetIncomingMessages()
	//go s.GenAndSendMailTest()
	go s.Send()
	for {
		select {
		case <-s.quitChan:
			log.Log().L("Data in quitChan. SMTP daemon not stopping...")
			//return
		}
	}
}

func (s *Server) RunSmtpServer() {
	if err := s.smtpServer.ListenAndServe(); err != nil {
		log.Log().E("Error in smtpServer.ListenAndServe: ", err)
		s.quitChan <- struct{}{}
	}
}

func (s *Server) GetIncomingMessages() {
	for data := range s.resultChannel {
		if data.Error != nil {
			log.Log().W("Wrong email got: ", data.Error)
			continue
		}
		log.Log().L("Got a message")
		log.Log().L(data.Email.Stringify())
		s.incomingQueueChan.In <- data.Email
	}
}

// PrintAndForward : gets message from MailSender, prints it and resends to IncomingQueue
func (s *Server) Send() {
	for pack := range s.mailSenderChan.Out {
		email := pack.(post.Email)
		//s.incomingQueueChan.In <- email
		fmt.Println("GOing to send a message")
		sender := NewSMTPSender(post.Conf.Login, post.Conf.Password, post.Conf.Host, post.Conf.Port, s.resultChannel)
		err := sender.Send(email.From, []string{email.To}, email.Subject, []byte(email.Body))
		if err != nil {
			log.Log().E("Cannot send email: ", err)
			s.getAndSendErrorMessage(err, email)
		} else {
			log.Log().L("Email sent")
		}
	}
}

func (s *Server) getAndSendErrorMessage(err error, message post.Email) error{
	message.To = message.From
	message.From = "mailder-daemon@nl-mail.ru"
	newBody := "Sorry, but we cannot delivery your message since:\n " + err.Error()
	newBody += "\n------Message-----\n"
	//newBody += message.Body
	new := gomail.NewMessage()
	//new.SetHeader("From", from.From, name)
	new.SetHeader("From", message.From)
	new.SetHeader("To", message.To)
	new.SetHeader("Subject", message.Subject)
	new.SetBody("text/plain", newBody)
	var bodyWriter bytes.Buffer
	new.WriteTo(&bodyWriter)
	message.Body = bodyWriter.String()
	log.Log().L("Created an error message: ", newBody)
	log.Log().L("New body end")
	s.mailSenderChan.Out <- message
	log.Log().E("Sent a message about error")
	return nil
}

func (s *Server) GenAndSendMailTest() {
	for {
		email := post.Email{From:"ivan", To:"ian@mam.sas", Body:`DKIM-Signature: v=1; a=rsa-sha256; q=dns/txt; c=relaxed/relaxed; d=mail.ru; s=mail2;
		h=ReSent-From:ReSent-To:Content-Type:Resent-Message-ID:Message-ID:Reply-To:Resent-Date:Date:MIME-Version:Subject:To:From; bh=2OIh9hyye6Hw4lQMPR3Loenu9B/A2RdXlp2hZL8j/Uw=;
		b=Nqm4jK9vGiLpnjXuTbcctSZNA0g3rv1SXCdlcwiFvFMb1H88/D1EC0GePS3lH+weTtlLpS+v3T87BZIbsMDAL54Kv66hR0SOpS6OVsiPwH674ERSc+3jUWLor3majzZqMIwcOFJ4SvAjHkbA0r3IdMCZEIve2VOlnSEjLeQ/wbk=;
	Received: by f100.i.mail.ru with local (envelope-from <ivanov.vanya.111@mail.ru>)
		id 1iRMOU-0004Ra-II
		for aa@nlmail.ddns.net; Sun, 03 Nov 2019 23:22:38 +0300
	Delivered-To: ivanov.vanya.111@mail.ru
	X-Received: by f147.i.mail.ru with local (envelope-from <ivanov.vanya.111@mail.ru>)
		id 1iRMMG-0007ME-Qb
		for ivanov.vanya.111@mail.ru; Sun, 03 Nov 2019 23:20:20 +0300
	X-Received: by e.mail.ru with HTTP;
		Sun, 03 Nov 2019 23:20:20 +0300
	From: =?UTF-8?B?0LjQstCw0L0g0LjQstCw0L3QvtCy?= <ivanov.vanya.111@mail.ru>
	X-Original-From: =?UTF-8?B?0LjQstCw0L0g0LjQstCw0L3QvtCy?= <ivanov.vanya.111@mail.ru>
	To: =?UTF-8?B?YWE=?= <aa@nlmail.ddns.net>
	Subject: =?UTF-8?B?VGVzdFN1YmplY3Q=?=
	MIME-Version: 1.0
	X-Mailer: Mail.Ru Mailer 1.0
	Date: Sun, 03 Nov 2019 23:22:38 +0300
	Resent-Date: Sun, 03 Nov 2019 23:20:20 +0300
	Reply-To: =?UTF-8?B?0LjQstCw0L0g0LjQstCw0L3QvtCy?= <ivanov.vanya.111@mail.ru>
	X-Priority: 3 (Normal)
	Message-ID: <1572812558.674382473@f100.i.mail.ru>
	Resent-Message-ID: <1572812420.76517441@f147.i.mail.ru>
	Content-Type: multipart/mixed;
		boundary="----e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM-WvXB9DzmKPIKmbn9-1572812420"
	X-Mailru-Intl-Transport: d,03c87f6
	Received: by e.mail.ru with HTTP;
		Sun, 03 Nov 2019 23:22:38 +0300
	ReSent-To: =?UTF-8?B?YWE=?= <aa@nlmail.ddns.net>
	ReSent-From: =?UTF-8?B?0LjQstCw0L0g0LjQstCw0L3QvtCy?= <ivanov.vanya.111@mail.ru>
	Authentication-Results: f100.i.mail.ru; auth=pass smtp.auth=ivanov.vanya.111@mail.ru smtp.mailfrom=ivanov.vanya.111@mail.ru
	X-77F55803: 82CEFCC3E280D3DB7F9F52485CB584D7EA5C12FFEB6BB2620A9D9B7F6070941CA011FE5DB9A1B056C5F29F942716E80B660823B2EEFD31DC
	X-7FA49CB5: 70AAF3C13DB7016878DA827A17800CE7D3037527CC315E3AD82A6BABE6F325ACA01ED31736435A1FBFD28B28ED4578739E625A9149C048EEFAD5A440E159F97D29508FF2E8683A3EB287FD4696A6DC2FA8DF7F3B2552694A4E2F5AFA99E116B42401471946AA11AF23F8577A6DFFEA7C5C0AD7D016C066E38F08D7030A58E5AD6BA297DBC24807EAA9D420A4CFB5DD3E02BCEA42C61739E62DF47820CA8ADC4182261AA87FD8CD6F8941B15DA834481FA18204E546F3947CD2DCF9CF1F528DBCF6B57BC7E64490618DEB871D839B7333395957E7521B51C2545D4CF71C94A83E9FA2833FD35BB23D27C277FBC8AE2E8B974A882099E279BDA471835C12D1D977C4224003CC836476C0CAF46E325F83A522CA9DD8327EE4931B544F03EFBC4D57EDE37B0041641260AD99DB1B8270F1D1731C566533BA786A40A5AABA2AD371193C9F3DD0FB1AF5EBE26B79914A659FB62623479134186CDE6BA297DBC24807EABDAD6C7F3747799A
	X-Mailru-MI: 800
	X-Mras: OK
	X-Spam: undefined
	
	
	------e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM-WvXB9DzmKPIKmbn9-1572812420
	Content-Type: multipart/alternative;
		boundary="--ALT--e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM1572812420"
	
	
	----ALT--e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM1572812420
	Content-Type: text/plain; charset=utf-8
	Content-Transfer-Encoding: base64
	
	SGVsbG8KCgotLSAK0JjQstCw0L0g0JrQvtGH0YPQsdC10Lk=
	----ALT--e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM1572812420
	Content-Type: text/html; charset=utf-8
	Content-Transfer-Encoding: base64
	
	CjxIVE1MPjxCT0RZPkhlbGxvPGJyPjxicj48YnI+LS0gPGJyPtCY0LLQsNC9INCa0L7Rh9GD0LHQ
	tdC5PC9CT0RZPjwvSFRNTD4K
	----ALT--e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM1572812420--
	
	------e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM-WvXB9DzmKPIKmbn9-1572812420
	Content-Type: application/octet-stream; name="=?UTF-8?B?aWdvci5t?="
	Content-Disposition: attachment; filename="=?UTF-8?B?aWdvci5t?="
	Content-Transfer-Encoding: base64
	
	V3A9dGYoWzM2MF0sWzUwLDBdKTsNCld6PXRmKFsxXSxbMC4wMDAyOCwwLjEwMjgsMV0pOw0KV289
	dGYoWzM2MF0sWzAuMTQsNi40LDUwLDBdKTsNCmJvZGUoV3AsV3osV28pOw0KbGVnZW5kKCfLwNTX
	1SDo7fLl4/Dg8u7w4CcsJ8vA1NfVIOTi7unt7uPuIODv5fDo7uTo9+Xx6u7j7iDn4uXt4CcsJ8vA
	1NfVIOLx5ekg8ejx8uXs+ycpOw0KdGl0bGUoJ8vA1NfVIPDg7Ort8/Lu6SDx6PHy5ez7Jyk7DQpn
	cmlkIG9uOw0KZ3JpZCBtaW5vcjs=
	------e10z6wlYbnF4v83EHs1FDDUXlsHUW4EM-WvXB9DzmKPIKmbn9-1572812420--`}
		s.incomingQueueChan.In <- email
		time.Sleep(5000 * time.Millisecond)
	}
}