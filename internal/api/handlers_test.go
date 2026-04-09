package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSubscribe_ValidationErrors перевіряє, як API реагує на неповні дані
func TestSubscribe_ValidationErrors(t *testing.T) {
	// Створюємо порожній обробник (без реальної бази і клієнта GitHub),
	// оскільки логіка перевірки на порожні поля спрацьовує ДО звернення до БД.
	handler := &Handler{}

	// Це наша "таблиця" тестових сценаріїв
	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
	}{
		{
			name:           "Empty email and repo",
			payload:        map[string]string{"email": "", "repository": ""},
			expectedStatus: http.StatusBadRequest, // Очікуємо статус 400
		},
		{
			name:           "Missing email",
			payload:        map[string]string{"email": "", "repository": "golang/go"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing repository",
			payload:        map[string]string{"email": "test@example.com", "repository": ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		// t.Run запускає кожен сценарій як окремий підтест
		t.Run(tt.name, func(t *testing.T) {
			// Перетворюємо тестові дані в JSON
			body, _ := json.Marshal(tt.payload)

			// Створюємо фейковий запит
			req, _ := http.NewRequest(http.MethodPost, "/api/subscribe", bytes.NewBuffer(body))

			// Створюємо фейковий "записувач" відповіді сервера
			rr := httptest.NewRecorder()

			// Викликаємо нашу функцію
			handler.Subscribe(rr, req)

			// Перевіряємо, чи отримали ми очікуваний статус
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Обробник повернув неправильний статус: отримали %v, очікували %v", status, tt.expectedStatus)
			}
		})
	}
}
