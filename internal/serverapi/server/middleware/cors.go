package middleware

import (
	"2019_2_Next_Level/internal/serverapi/config"
	"2019_2_Next_Level/internal/serverapi/log"
	"net/http"

	"github.com/gorilla/mux"
)

func CorsMethodMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			origin := r.Header.Get("Origin")
			log.Log().L("origin: ", origin)
			if !config.Conf.HttpConfig.Whitelist[origin] {
				log.Log().I("Not in whitelist: ", origin)
				http.Error(w, "Not in whitelist", http.StatusForbidden)
				return
			}
			headers.Add("Access-Control-Allow-Origin", origin)
			headers.Add("Access-Control-Allow-Credentials", "true")
			headers.Add("Access-Control-Allow-Headers", "Content-Type")
			headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

			next.ServeHTTP(w, r)
		})
	}
}
