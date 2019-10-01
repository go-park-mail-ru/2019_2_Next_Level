package daemon

import (
	"encoding/json"
	"net/http"
)

const (
	ErrorNoPermission = "No permission"
	ErrorInternal     = "Internal error"
	ErrorBadRequest   = "Bad request"
)

type Error struct {
	Value string `json:"error"`
}

func (e *Error) Send(w *http.ResponseWriter) {
	str, _ := json.Marshal(e)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusBadRequest)
	(*w).Write(str)
}
