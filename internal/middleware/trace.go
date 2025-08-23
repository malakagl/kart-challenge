package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/malakagl/kart-challenge/pkg/constants"
)

// Trace ensures every request has a trace ID.
func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("x-request-id")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		w.Header().Set("x-request-id", traceID)
		ctx := context.WithValue(r.Context(), constants.TraceIDKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
