package smtpd

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/smtpd/worker"
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

type SMTPSender struct {
	login string
	password string
	host string
	port string
	localOutChan chan worker.EmailNil
}

func NewSMTPSender(login string, password string, host string, port string, localOutChan chan worker.EmailNil) *SMTPSender {
	return &SMTPSender{login: login, password: password, host: host, port: port, localOutChan:localOutChan}
}

func (s *SMTPSender) Send(from string, to []string, subject string, body []byte) error {
	//auth := smtp.PlainAuth("", s.login, s.password, s.host)
	//err := smtp.SendMail(s.host + ":" + s.port, auth, from, to, body)

	//msg := composeMimeMail(to[0], from, subject, string(body))
	msg := body

	mx, err := getMXRecord(to[0])
	fmt.Println("MX: ", mx)
	if err != nil {
		if err.Error()=="Local" {
			s.localOutChan <- worker.EmailNil{
				Email: post.Email{
					From:    from,
					To:      to[0],
					Body:    string(body),
					Subject: subject,
				},
				Error: nil,
			}
			fmt.Println("Sent local")
			return nil;
		}
		return err
	}

	err = smtp.SendMail(mx+":25", nil, from, to, msg)
	if err != nil {
		return err;
	}
	return err
}

func getMXRecord(to string) (mx string, err error) {
	cache := map[string]string {
		"mail.nl-mail.ru": "mail.nl-mail.ru",
		"nl-mail.ru": "mail.nl-mail.ru",
		"nlmail.hldns.ru" : "nlmail.hldns.ru",
	}
	var e *mail.Address
	e, err = mail.ParseAddress(to)
	if err != nil {
		return
	}

	domain := strings.Split(e.Address, "@")[1]

	var mxs []*net.MX
	//mxs, err = net.LookupMX(domain)
	fmt.Println("Domain: ", domain)
	if _, ok := cache[domain]; ok {
		return "",fmt.Errorf("Local")
	}
	mxs, err = LookupMXCached(domain, cache)

	if err != nil || len(mxs) == 0 {
		fmt.Println("Error 1: ", err)
		return
	}
	host := mxs[0].Host
	if host[len(host)-1] == '.' {
		host = host[:len(host)-1]
	}
	mx = host

	return
}

// Never fails, tries to format the address if possible
func formatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>")
}

func composeMimeMail(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = formatEmailAddress(from)
	header["To"] = formatEmailAddress(to)
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}

func LookupMXCached(domain string, cache map[string]string) ([]*net.MX, error) {
	if addr, ok := cache[domain]; ok {
		return []*net.MX{
			&net.MX{Host:addr},
		}, nil
	}
	return net.LookupMX(domain)
}