package worker

import (
	"2019_2_Next_Level/internal/post"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
)

// A Session is returned after successful login.
// Satisfies the interface smtp.Session{ Mail(), Rcpt(), Data(), Reset(), Logout() }
type Session struct {
	ID     uuid.UUID
	email  post.Email
	atEnd  func(EmailNil)
	isDone bool
}

// Init : initialization of the Session
func (s Session) Init() Session {
	s.ID, _ = uuid.NewUUID()
	return s
}

// Mail : callback triggered when extern server sends command "MAIL FROM:"
func (s *Session) Mail(from string) error {
	log.Println("Mail from:", from)
	s.email.From = from
	return nil
}

// Rcpt : callback triggered when extern server sends command "rcpt to:"
func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	s.email.To = to
	return nil
}

// Data : callback triggered when extern server sends "." (end of body transmittion) after block "DATA"
func (s *Session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.email.Body = string(b)

	//	maybe we should send resulted message here?

	return nil
}

// Reset : discard currently processed message
func (s *Session) Reset() {
	fmt.Println("Reset")
	s.SendResults()
}

// Logout : callback triggered when extern server closes the connection
func (s *Session) Logout() error {
	fmt.Println("Logout")
	s.SendResults()
	return nil
}

// SendResults : checks if session is active and triggers result-returning callback
func (s *Session) SendResults() {
	if !s.isDone { // почему-то Logout() вызывается дважды для каждого сообщения
		s.atEnd(EmailNil{Email: s.email})
	}
	s.isDone = true
}
