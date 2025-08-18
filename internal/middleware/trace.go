package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/malakagl/kart-challenge/pkg/constants"
)

// TraceMiddleware ensures every request has a trace ID.
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", traceID)
		ctx := context.WithValue(r.Context(), constants.TraceIDKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
