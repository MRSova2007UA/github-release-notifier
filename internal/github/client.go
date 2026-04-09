package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client — це структура нашого клієнта для GitHub API
type Client struct {
	httpClient *http.Client
	token      string // Знадобиться для збільшення лімітів
}

// NewClient створює новий екземпляр клієнта
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Щоб запити не висіли вічно
		},
		token: token,
	}
}

// ValidateRepo перевіряє правильність формату (owner/repo) та існування репозиторію
func (c *Client) ValidateRepo(repoName string) (int, error) {
	// Перевіряємо формат: має бути рівно дві частини, розділені слешем
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return http.StatusBadRequest, fmt.Errorf("неправильний формат, очікується owner/repo")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Якщо є токен — додаємо його для авторизації (щоб не зловити ліміт 60 запитів)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	// Обробка помилок зовнішнього API (Вимога #7)
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
		return http.StatusTooManyRequests, fmt.Errorf("перевищено ліміт запитів до GitHub API")
	}

	if resp.StatusCode == http.StatusNotFound {
		return http.StatusNotFound, fmt.Errorf("репозиторій не знайдено на GitHub")
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("неочікувана помилка від GitHub: %d", resp.StatusCode)
	}

	return http.StatusOK, nil
}

// ReleaseData - структура для парсингу відповіді з релізом
type ReleaseData struct {
	TagName string `json:"tag_name"`
}

// GetLatestRelease дістає останній тег релізу (наприклад, "v1.2.3")
func (c *Client) GetLatestRelease(repoName string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Репозиторій є, але релізів у ньому ще немає
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("помилка отримання релізу: %d", resp.StatusCode)
	}

	var release ReleaseData
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("помилка читання JSON: %v", err)
	}

	return release.TagName, nil
}
