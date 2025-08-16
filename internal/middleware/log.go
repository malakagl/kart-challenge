package middleware

import (
	"net/http"
	"time"

	logging "github.com/malakagl/kart-challenge/pkg/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			logging.Logger.Info().Msgf("Request %s %s %s processed in %s", r.Method, r.URL.Path, r.RemoteAddr, time.Since(startTime))
		}()

		next.ServeHTTP(w, r)
	})
}
