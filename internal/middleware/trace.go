package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Trace ensures every request has a trace ID and propagates context.
func Trace(next http.Handler) http.Handler {
	tracer := otel.Tracer("kart-challenge")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract any existing trace context from incoming headers
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// Start a new span for this HTTP request
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// Add standard HTTP attributes
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.Path),
			attribute.String("http.client_ip", r.RemoteAddr),
		)

		// Propagate trace ID in response headers for clients
		traceID := span.SpanContext().TraceID().String()
		w.Header().Set("x-request-id", traceID)
		w.Header().Set("traceparent", span.SpanContext().TraceID().String()) // optional for W3C compliance

		// Pass the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
