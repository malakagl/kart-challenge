package middleware

import (
	"net/http"

	"github.com/malakagl/kart-challenge/pkg/models/dto/response"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api_key")
		if apiKey == "" || !isValidApiKey(apiKey) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			response.Error(w, http.StatusUnauthorized, "AuthError", http.StatusText(http.StatusUnauthorized))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidApiKey(key string) bool {
	return key == "apitest" // Replace with actual API key validation logic
}
