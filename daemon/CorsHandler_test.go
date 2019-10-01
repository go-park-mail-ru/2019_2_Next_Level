package daemon

import (
	"net/http/httptest"
	"testing"
)

func TestPreflightHandler(t *testing.T) {
	r := httptest.NewRequest("POST", "/auth.signin", nil)
	w := httptest.NewRecorder()

	(&CorsHandler{}).preflightHandler(w, r)

	headers := w.Header()
	if headers.Get("Access-Control-Allow-Origin") != "http://localhost:3000" ||
		headers.Get("Access-Control-Allow-Credentials") != "true" ||
		headers.Get("Access-Control-Allow-Headers") != "Content-Type" ||
		headers.Get("Access-Control-Allow-Methods") != "GET, POST, OPTIONS" {
		t.Error("Wrong headers got")
	}
}
