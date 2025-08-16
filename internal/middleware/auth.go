package middleware

import "net/http"

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api_key")
		if apiKey == "" || !isValidApiKey(apiKey) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidApiKey(key string) bool {
	return key == "apitest" // Replace with actual API key validation logic
}
