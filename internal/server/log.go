package server

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			log.Printf("Request %s %s %s processed in %s\n", r.Method, r.URL.Path, r.RemoteAddr, time.Since(startTime))
		}()
		next.ServeHTTP(w, r)
	})
}
