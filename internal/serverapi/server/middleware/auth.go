package middleware

import (
	"2019_2_Next_Level/internal/post/log"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	hr "2019_2_Next_Level/internal/serverapi/server/Error/httpError"
	"2019_2_Next_Level/pkg/HttpTools"
	"net/http"

	"github.com/gorilla/mux"
)

func AuthentificationMiddleware(authCase auth.Usecase) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session-id")
			if err != nil {
				(&HttpTools.Response{}).SetWriter(w).SetError(hr.GetError(hr.BadSession)).Send()
				return
			}
			login, res := authCase.CheckAuth(cookie.Value)
			if res != nil {
				(&HttpTools.Response{}).SetWriter(w).SetError(hr.GetError(hr.BadSession)).Send()
				log.Log().I("No permission")
				return
			}
			r.Header.Set("X-Login", login)
			next.ServeHTTP(w, r)
		})

	}
}
