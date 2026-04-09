package api

import (
	"encoding/json"
	"net/http"

	"github-release-notifier/internal/github"
	"github-release-notifier/internal/repository"
)

// Handler містить посилання на БД та GitHub клієнт
type Handler struct {
	repo     *repository.Repository
	ghClient *github.Client
}

// NewHandler створює новий екземпляр обробника
func NewHandler(repo *repository.Repository, ghClient *github.Client) *Handler {
	return &Handler{
		repo:     repo,
		ghClient: ghClient,
	}
}

// SubscribeRequest - це структура JSON, який ми будемо отримувати від користувача
type SubscribeRequest struct {
	Email      string `json:"email"`
	Repository string `json:"repository"`
}

// Subscribe - функція, яка обробляє запит POST /subscribe
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неправильний формат JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Repository == "" {
		http.Error(w, `{"error": "Поля email та repository є обов'язковими"}`, http.StatusBadRequest)
		return
	}

	status, err := h.ghClient.ValidateRepo(req.Repository)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, status)
		return
	}

	latestTag, err := h.ghClient.GetLatestRelease(req.Repository)
	if err != nil {
		http.Error(w, `{"error": "Помилка отримання релізу з GitHub"}`, http.StatusInternalServerError)
		return
	}

	if err := h.repo.SubscribeUser(req.Email, req.Repository, latestTag); err != nil {
		http.Error(w, `{"error": "Помилка збереження в БД"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Успішно підписано на оновлення!"}`))
}
