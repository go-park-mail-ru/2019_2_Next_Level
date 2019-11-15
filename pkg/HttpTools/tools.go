package HttpTools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Inflate struct from json written in body
func StructFromBody(r http.Request, s interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, s)
}

func BodyFromStruct(w http.ResponseWriter, s interface{}) error {
	js, err := json.Marshal(s)
	if err != nil {
		return err
	}
	w.Write(js)
	w.Header().Add("Content-Type", "application/json")
	return nil
}

type Answer interface {
	SetStatus(string)
}

const (
	defaultStatus = "OK"
)

type Response struct {
	Error  interface{}
	writer http.ResponseWriter
}

func (r *Response) SetWriter(writer http.ResponseWriter) *Response {
	r.writer = writer
	return r
}
func (r *Response) SetError(err interface{}) *Response {
	r.Error = err
	return r
}

func (r *Response) Copy() Response {
	return *r
}

func (r *Response) Send() {
	if r.Error == nil {
		fmt.Println("Nil error")
	}
	body, err := json.Marshal(r.Error)
	if err != nil {
		fmt.Println("Cannot encode json")
		return
	}
	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.Write(body)
}

func (r *Response) String() string {
	body, _ := json.Marshal(r.Error)
	return string(body)
}
