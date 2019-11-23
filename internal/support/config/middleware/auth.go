package middleware

//func AuthentificationMiddleware(authCase auth.Usecase) mux.MiddlewareFunc {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			cookie, err := r.Cookie("session-id")
//			if err != nil {
//				(&HttpTools.Response{}).SetWriter(w).SetError(hr.GetError(hr.BadSession)).Send()
//				return
//			}
//			//login, res := authCase.CheckAuth(cookie.Value)
//			//if res != nil {
//			//	(&HttpTools.Response{}).SetWriter(w).SetError(hr.GetError(hr.BadSession)).Send()
//			//	log.Log().I("No permission")
//			//	return
//			//}
//			r.Header.Set("X-Login", login)
//			next.ServeHTTP(w, r)
//		})
//
//	}
//}
