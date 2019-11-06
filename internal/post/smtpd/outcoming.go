package smtpd

import (
	"net/smtp"
)

type SMTPSender struct {
	login string
	password string
	host string
	port string
}

func NewSMTPSender(login string, password string, host string, port string) *SMTPSender {
	return &SMTPSender{login: login, password: password, host: host, port: port}
}

func (s *SMTPSender) Send(from string, to []string, body []byte) error {
	auth := smtp.PlainAuth("", s.login, s.password, s.host)
	err := smtp.SendMail(s.host + ":" + s.port, auth, from, to, body)
	return err
}