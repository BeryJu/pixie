package server

import (
	log "github.com/sirupsen/logrus"

	"net/http"
)

func logging(logger *log.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.WithFields(log.Fields{
				"method": r.Method,
				"remote": r.RemoteAddr,
			}).Info(r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
