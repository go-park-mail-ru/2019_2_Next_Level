package post

import (
	"fmt"
)

type Email struct {
	From string
	To   string
	Body string
	Subject string
}

func (e *Email) Stringify() string {
	return fmt.Sprintf("From: %s\nTo: %s\nBody: %s", e.From, e.To, e.Body)
}
