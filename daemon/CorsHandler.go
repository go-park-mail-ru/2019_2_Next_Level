package daemon

import (
	"log"
	"net/http"
)

type CorsHandler struct {
}

func (h *CorsHandler) preflightHandler(w http.ResponseWriter, r *http.Request) {
	whitelist := map[string]bool{
		config.FrontendUrl:                              true,
		"https://next-level-mail.kerimovdev.now.sh/":    true,
		"https://next-level-mail.ivanovvanya111.now.sh": true,
	}
	headers := w.Header()
	// headers.Add("Access-Control-Allow-Origin", config.FrontendUrl)
	origin := r.Header.Get("Origin")
	if origin == "" {
		log.Println("No Origin header")
		return
	}
	log.Println("origin: ", origin)
	if whitelist[origin] {
		headers.Add("Access-Control-Allow-Origin", origin)
		log.Println("In whitelist")
	} else {
		log.Println("Not in whitelist: ", origin)
	}
	headers.Add("Access-Control-Allow-Credentials", "true")
	headers.Add("Access-Control-Allow-Headers", "Content-Type")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}
