package middleware

import (
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// type MiddlewareFunc func(http.Handler) http.Handler

func AuthentificationMiddleware(authCase auth.Usecase) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie := http.Cookie{}

			res := authCase.CheckAuthorization(&cookie)
			if res != nil {
				fmt.Println("No permission")
				return
			}
			next.ServeHTTP(w, r)
		})

	}
}

// func AuthentificationMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		auth.Usecase
// 	})
// }
