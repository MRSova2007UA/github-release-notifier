package api

import (
	"encoding/json"
	"github-release-notifier/internal/metrics"
	"net/http"
	"strings"

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

	// Розбиваємо рядок "owner/repo" на дві частини
	parts := strings.Split(req.Repository, "/")
	if len(parts) != 2 {
		http.Error(w, `{"error": "Неправильний формат репозиторію. Очікується owner/repo (наприклад, golang/go)"}`, http.StatusBadRequest)
		return
	}
	owner := parts[0]
	repoName := parts[1]

	// Перевіряємо, чи існує репозиторій (з використанням Redis кешу та контексту)
	exists, err := h.ghClient.CheckRepoExists(r.Context(), owner, repoName)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, `{"error": "Репозиторій не знайдено"}`, http.StatusNotFound)
		return
	}

	// Отримуємо останній реліз
	latestTag, err := h.ghClient.GetLatestRelease(r.Context(), owner, repoName)
	if err != nil {
		http.Error(w, `{"error": "Помилка отримання релізу з GitHub"}`, http.StatusInternalServerError)
		return
	}

	// Зберігаємо в базу даних
	if err := h.repo.SubscribeUser(req.Email, req.Repository, latestTag); err != nil {
		http.Error(w, `{"error": "Помилка збереження в БД"}`, http.StatusInternalServerError)
		return
	}

	// Рахуємо успішну підписку в Prometheus (зверни увагу, слешів тут немає)
	metrics.TotalSubscriptions.Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Успішно підписано на оновлення!"}`))
}
