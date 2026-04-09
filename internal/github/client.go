package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github-release-notifier/internal/metrics" // ДОДАНО
	"github.com/redis/go-redis/v9"
)

// Client - структура клієнта для роботи з GitHub API та Redis
type Client struct {
	token      string
	httpClient *http.Client
	rdb        *redis.Client
}

// NewClient створює новий екземпляр клієнта
func NewClient(token string, rdb *redis.Client) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		rdb: rdb,
	}
}

// CheckRepoExists перевіряє чи існує репозиторій (з використанням Redis кешу)
func (c *Client) CheckRepoExists(ctx context.Context, owner, repo string) (bool, error) {
	// 1. Формуємо унікальний ключ для кешу
	cacheKey := fmt.Sprintf("repo_exists:%s/%s", owner, repo)

	// 2. Шукаємо в Redis
	val, err := c.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Printf("Cache HIT! Репозиторій %s/%s знайдено в пам'яті", owner, repo)
		return val == "true", nil
	}

	// 3. Якщо в кеші немає — ідемо в GitHub
	log.Printf("Cache MISS. Робимо запит в GitHub для %s/%s...", owner, repo)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		metrics.GitHubAPICalls.WithLabelValues("error").Inc() // ДОДАНО
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		metrics.GitHubAPICalls.WithLabelValues("rate_limited").Inc() // ДОДАНО
		return false, fmt.Errorf("досягнуто ліміт запитів GitHub API")
	}

	metrics.GitHubAPICalls.WithLabelValues("success").Inc() // ДОДАНО
	exists := resp.StatusCode == http.StatusOK

	// 4. Записуємо результат в Redis на 10 хвилин
	c.rdb.Set(ctx, cacheKey, fmt.Sprintf("%t", exists), 10*time.Minute)

	return exists, nil
}

// Release - структура для парсингу відповіді з релізом
type Release struct {
	TagName string `json:"tag_name"`
}

// GetLatestRelease отримує останній реліз репозиторію
func (c *Client) GetLatestRelease(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		metrics.GitHubAPICalls.WithLabelValues("error").Inc() // ДОДАНО (тут теж корисно рахувати)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		metrics.GitHubAPICalls.WithLabelValues("success").Inc() // ДОДАНО
		return "", nil                                          // Релізів ще немає
	}

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		metrics.GitHubAPICalls.WithLabelValues("rate_limited").Inc() // ДОДАНО
		return "", fmt.Errorf("досягнуто ліміт запитів GitHub API")
	}

	if resp.StatusCode != http.StatusOK {
		metrics.GitHubAPICalls.WithLabelValues("error").Inc() // ДОДАНО
		return "", fmt.Errorf("неочікуваний статус від GitHub: %d", resp.StatusCode)
	}

	metrics.GitHubAPICalls.WithLabelValues("success").Inc() // ДОДАНО

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}

	return rel.TagName, nil
}
