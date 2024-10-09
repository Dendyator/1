package internalhttp

import (
	"net/http"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
)

func loggingMiddleware(logg *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logg.Infof("%s [%s] %s %s %s %d %s",
				r.RemoteAddr,
				start.Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.URL.Path,
				r.Proto,
				http.StatusOK,
				time.Since(start),
			)
		})
	}
}
