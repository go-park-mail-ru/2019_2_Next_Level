package middleware

import (
	"2019_2_Next_Level/internal/serverapi/log"
	"2019_2_Next_Level/internal/serverapi/server/metrics"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	hr "2019_2_Next_Level/internal/serverapi/server/HttpError"
	"github.com/gorilla/mux"
)
type statusWriter struct {
	http.ResponseWriter
	status int
	length int
	body []byte
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.body = b
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func AccessLogMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := r.URL
			log.Log().I(
				fmt.Sprintf(
					"%s Request\nTo: %s\nOrigin: %s\nMethod: %s\nAgent: %s\n",
						time.Now().String(),
						url,
						r.Header.Get("Origin"),
						r.Method,
						r.Header.Get("User-Agent")))
			sw := &statusWriter{}
			sw.ResponseWriter = w
			//starttime := time.Now()
			next.ServeHTTP(sw, r)
			var res hr.HttpResponse
			_ = json.Unmarshal(sw.body, &res)
			metrics.Hits.WithLabelValues(res.Status, r.URL.Path).Inc()
		})

	}
}
