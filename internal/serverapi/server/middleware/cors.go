package middleware

import (
	"2019_2_Next_Level/internal/serverapi/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func CorsMethodMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			origin := r.Header.Get("Origin")
			if origin == "" {
				log.Println("No Origin header")
				return
			}
			log.Println("origin: ", origin)
			dd := config.Conf.HttpConfig.Whitelist
			fmt.Println(dd)
			if !config.Conf.HttpConfig.Whitelist[origin] {
				log.Println("Not in whitelist: ", origin)
				http.Error(w, "Not in whitelist", http.StatusForbidden)
				return
			}
			log.Println("In whitelist")
			headers.Add("Access-Control-Allow-Origin", origin)
			headers.Add("Access-Control-Allow-Credentials", "true")
			headers.Add("Access-Control-Allow-Headers", "Content-Type")
			headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

			next.ServeHTTP(w, r)
		})
	}
}
