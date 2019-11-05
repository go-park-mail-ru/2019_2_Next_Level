package worker

import (
	"fmt"
	"io"

	"github.com/emersion/go-smtp"
)

type MockWorker struct{}

func (w *MockWorker) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	// if username != "username@nextmail.ru" || password != "password" {
	// 	return nil, errors.New("Invalid username or password")
	// }
	session := MockSession{}

	return &session, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (w *MockWorker) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	fmt.Println("AnonimousLogin")
	session := MockSession{}
	return &session, nil
}

type MockSession struct {
	resChan chan string
}

func (s MockSession) Init() MockSession {
	s.resChan <- "Init"
	return s
}

func (s *MockSession) Mail(from string) error {
	s.resChan <- "Mail"
	return nil
}

func (s *MockSession) Rcpt(to string) error {
	s.resChan <- "Rcpt"
	return nil
}

func (s *MockSession) Data(r io.Reader) error {
	s.resChan <- "Data"
	return nil
}

func (s *MockSession) Reset() {
	s.resChan <- "Reset"
}

func (s *MockSession) Logout() error {
	s.resChan <- "Logout"
	return nil
}

func (s *MockSession) SendResults() {
	s.resChan <- "SendResults"
}
