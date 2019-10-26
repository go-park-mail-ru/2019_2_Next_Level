package main

import (
	"encoding/json"
	"fmt"
)

func Construct(status string, args map[string]interface{}) {
	args["status"] = status
	js, _ := json.Marshal(args)
	fmt.Println(string(js))
}

type T struct {
	Data int
	Name string
}

func (t *T) Set(a int) {
	t.Data = a
}

func Parse(data []byte, str interface{}) error {
	return json.Unmarshal(data, str)
}

func main() {
	// js, _ := json.Marshal("fullName", "password")
	// fmt.Println(string(js))
	str := `{"data":12, "name": "ian"}`
	t := T{}
	// json.Unmarshal([]byte(str), &t)
	Parse([]byte(str), &t)
	fmt.Println(t)

}
