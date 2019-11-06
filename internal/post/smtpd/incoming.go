package smtpd

import (
	"github.com/emersion/go-smtp"
	"time"
)

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
