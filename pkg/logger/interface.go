package logger

type ILog interface {
	SetPrefix(string)
	Println(...interface{})
	I(...interface{})
	E(...interface{})
	W(...interface{})
	L(...interface{})
}
