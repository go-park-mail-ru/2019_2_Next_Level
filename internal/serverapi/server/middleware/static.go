package middleware

import (
	"net/http"
)

func StaticMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path[len(path)-1] == '/'{
				path = path + "index.html"
			}
			r.URL.Path = path
			next.ServeHTTP(w, r)
		})

	}
}
