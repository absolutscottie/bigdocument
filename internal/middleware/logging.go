package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// LoggingMiddleware configures the provided http.Handler to log inforrmative
// messages when http requests are received.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			},
		).Info()
		next.ServeHTTP(w, r)
	})
}
