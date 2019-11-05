package worker

import (
	"2019_2_Next_Level/internal/post"
	"github.com/emersion/go-smtp"
)

// EmailNil : stores result of the receiving message:
// mail (inflated or empty) and error
type EmailNil struct {
	Email post.Email
	Error error
}

// Worker : The Backend implements SMTP server methods.
type Worker struct {
	session    Session
	OutChannel chan EmailNil
}

// Login : handles a login command with username and password.
func (w *Worker) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	// if username != "username@nextmail.ru" || password != "password" {
	// 	return nil, errors.New("Invalid username or password")
	// }
	w.session = Session{atEnd: w.ObtainResult}.Init()
	return &w.session, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (w *Worker) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	w.session = Session{atEnd: w.ObtainResult}.Init()
	return &w.session, nil
}

// ObtainResult : callback for Session to get a result of her lifecycle
func (w *Worker) ObtainResult(data EmailNil) {
	w.OutChannel <- data
}
