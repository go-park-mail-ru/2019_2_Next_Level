package logger

import (
	"fmt"
)


type Log struct {
	prefix   string
	template string
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) SetPrefix(prefix string) {
	l.prefix = prefix
	l.template = "%s: (%s) %s"
}

func (l *Log) L(a ...interface{}) {
	l.Println(a...)
}

func (l *Log) I(a ...interface{}) {
	l.Println("(info) " + fmt.Sprintln(a...))
}

func (l *Log) E(a ...interface{}) {
	l.Println( "(error) " + fmt.Sprintln(a...))
}

func (l *Log) W(a ...interface{}) {
	l.Println( "(warning) " + fmt.Sprintln(a...))
}

func (l *Log) Println(a ...interface{}) {
	fmt.Println(l.prefix + ": " + fmt.Sprintln(a...))
}
