package api

import (
	"net/http"
	"os"
)

// APIKeyMiddleware перевіряє наявність секретного ключа в заголовку запиту
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Читаємо правильний ключ з налаштувань (або беремо стандартний для тестів)
		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			expectedKey = "super-secret-key-123"
		}

		// Перевіряємо, чи передав клієнт ключ у заголовку X-API-Key
		providedKey := r.Header.Get("X-API-Key")

		if providedKey != expectedKey {
			// Якщо ключ неправильний, відбиваємо запит
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Доступ заборонено. Неправильний API ключ."}`))
			return
		}

		// Якщо все ок - пропускаємо запит далі до нашого Handler'а
		next.ServeHTTP(w, r)
	})
}
