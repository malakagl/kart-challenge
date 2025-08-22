package middleware

import (
	"net/http"
	"time"

	"github.com/malakagl/kart-challenge/pkg/log"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			log.WithCtx(r.Context()).Info().Msgf("Request %s %s %s processed in %s", r.Method, r.URL.Path, r.RemoteAddr, time.Since(startTime))
		}()

		next.ServeHTTP(w, r)
	})
}
