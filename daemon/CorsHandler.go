package daemon

import (
	"log"
	"net/http"
)

type CorsHandler struct {
}

func (h *CorsHandler) preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")
	if origin == "" {
		log.Println("No Origin header")
		return
	}
	log.Println("origin: ", origin)
	if !configuration.Whitelist[origin] {
		log.Println("Not in whitelist: ", origin)
		http.Error(w, "Not in whitelist", http.StatusForbidden)
		return
	}
	log.Println("In whitelist")
	headers.Add("Access-Control-Allow-Origin", origin)
	headers.Add("Access-Control-Allow-Credentials", "true")
	headers.Add("Access-Control-Allow-Headers", "Content-Type")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}
