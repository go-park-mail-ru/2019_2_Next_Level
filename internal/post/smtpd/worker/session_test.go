package worker

import (
	"strings"
	"testing"
)

func TestSeccion(t *testing.T) {
	session := Session{}
	session.atEnd = func(e EmailNil) {

	}
	session.Init()
	session.Mail("from")
	session.Rcpt("to")
	r := strings.NewReader("")
	session.Data(r)
	session.Reset()
	session.SendResults()
}
