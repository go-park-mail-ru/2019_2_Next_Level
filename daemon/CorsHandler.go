package daemon

import "net/http"

type CorsHandler struct {
}

func (h *CorsHandler) preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Access-Control-Allow-Origin", config.FrontendUrl)
	headers.Add("Access-Control-Allow-Credentials", "true")
	headers.Add("Access-Control-Allow-Headers", "Content-Type")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}
