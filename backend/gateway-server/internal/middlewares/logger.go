package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func RequestLogging(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			entry := log.WithFields(logrus.Fields{
				"method":   r.Method,
				"path":     r.URL.Path,
				"status":   rw.statusCode,
				"duration": duration.String(),
				"ip":       r.RemoteAddr,
			})

			if location := rw.Header().Get("Location"); location != "" {
				entry = entry.WithField("redirect_to", location)
			}

			if rw.statusCode >= 500 {
				entry.Error("server error")
			} else if rw.statusCode >= 400 {
				entry.Warn("client error")
			} else if rw.statusCode >= 300 {
				entry.Info("redirect")
			}
		})
	}
}
