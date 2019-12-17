package HttpTools

import (
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"net/http"
)

type Jsonisable interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
	MarshalEasyJSON(w *jwriter.Writer)
	UnmarshalEasyJSON(w *jlexer.Lexer)
}

type Response struct {
	body  Jsonisable
	status int
	writer http.ResponseWriter
}

func NewResponse(writer http.ResponseWriter) *Response {
	return &Response{writer: writer}
}

func (r *Response) SetWriter(writer http.ResponseWriter) *Response {
	r.writer = writer
	return r
}

func (r *Response) SetStatus(status int) *Response {
	r.status = status
	return r
}

func (r *Response) SetError(err Jsonisable) *Response {
	r.body = err
	return r
}

func (r *Response) SetContent(c Jsonisable) *Response {
	r.body = c
	return r
}

func (r *Response) Copy() Response {
	return *r
}

func (r *Response) Send() {
	if r.body == nil {
		fmt.Println("Nil error")
	}
	body, err := easyjson.Marshal(r.body)
	if err != nil {
		fmt.Println("Cannot encode json")
		return
	}
	r.writer.Header().Set("Content-Type", "application/json")
	if r.status == 0 {
		r.status = http.StatusOK
	}
	r.writer.WriteHeader(r.status)
	r.writer.Write(body)
}

func (r *Response) String() string {
	body, _ := json.Marshal(r.body)
	return string(body)
}

type Answer interface {
	SetStatus(string)
}

const (
	defaultStatus = "OK"
)