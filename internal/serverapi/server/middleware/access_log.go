package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func AccessLogMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := r.URL.String()
			fmt.Printf("%s Request on %s\n", time.Now().String(), url)
			next.ServeHTTP(w, r)
		})

	}
}
