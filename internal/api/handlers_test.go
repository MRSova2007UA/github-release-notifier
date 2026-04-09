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
	handler := &Handler{}

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
	}{
		{
			name:           "Empty email and repo",
			payload:        map[string]string{"email": "", "repository": ""},
			expectedStatus: http.StatusBadRequest,
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
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)

			req, _ := http.NewRequest(http.MethodPost, "/api/subscribe", bytes.NewBuffer(body))

			rr := httptest.NewRecorder()

			handler.Subscribe(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Обробник повернув неправильний статус: отримали %v, очікували %v", status, tt.expectedStatus)
			}
		})
	}
}
