package model

import (
	"2019_2_Next_Level/internal/post"
	"fmt"
	"strings"
	"time"
)

type Email struct {
	From string
	To   string
	Body string
	Header struct {
		From string
		To   []string
		Subject string
		ReplyTo []string
		WhenReceived time.Time
	}
	Id int
	IsRead bool
}

func (e *Email) SetFrom(from string) {
	e.From = from
}

func (e *Email) Stringify() string {
	return fmt.Sprintf("From: %s\nTo: %s\nBody: %s", e.From, e.To, e.Body)
}

func (e *Email) FromPostEmail(orig post.Email) Email {
	e.Body = orig.Body
	e.From = orig.From
	e.To = orig.To
	return *e
}

func (e *Email) ToPostEmail() post.Email {
	email := post.Email{
		From: e.From,
		To:   e.To,
		Body: e.Body,
	}
	return email
}

func (e *Email) Split(src string) (string, string) {
	splitter := strings.Index(src, "@")
	if splitter==-1{
		return src, ""
	}
	login := src[:splitter]
	host := src[splitter+1:]
	return login, host
}