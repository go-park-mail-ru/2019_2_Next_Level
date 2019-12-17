package HttpTools

import (
	"encoding/json"
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

// Inflate body with struct content
func BodyFromStruct(w http.ResponseWriter, s interface{}) error {
	js, err := json.Marshal(s)
	if err != nil {
		return err
	}
	w.Write(js)	//sets unchangeable status 200 if not set already
	return nil
}

