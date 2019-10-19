package logger

import (
	"fmt"
	"log"
)

// log "github.com/ivahaev/go-logger"

type Log struct {
	prefix   string
	template string
}

func (l *Log) SetPrefix(prefix string) {
	l.prefix = prefix
	l.template = "%s: (%s) %s"
	// log.SetPrefix("")
}

// D : debug level
func (l *Log) Debugf(template string, a ...interface{}) {
	// l.Printf(
	// 	"(debug) "+template,
	// 	a...,
	// )
	// log.Debugf(l.prefix+": "+template, a...)
}

func (l *Log) I(template string, a ...interface{}) {
	l.Printf(
		"(info) "+template,
		a...,
	)
}

func (l *Log) E(template string, a ...interface{}) {
	l.Printf(
		"(error) "+template,
		a...,
	)
}

func (l *Log) Printf(template string, a ...interface{}) {
	// log.Printf(
	// 	l.prefix+": "+template,
	// 	a...,
	// )
}

func (l *Log) Println(a ...interface{}) {
	log.Println(l.prefix + ": " + fmt.Sprintln(a...))
}
