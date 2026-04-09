package api

import (
	"net/http"
	"os"
)

// APIKeyMiddleware перевіряє наявність секретного ключа в заголовку запиту
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			expectedKey = "super-secret-key-123"
		}

		providedKey := r.Header.Get("X-API-Key")

		if providedKey != expectedKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Доступ заборонено. Неправильний API ключ."}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
