package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Лічильник загальної кількості підписок
	TotalSubscriptions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "github_notifier_subscriptions_total",
		Help: "Загальна кількість успішних підписок на релізи",
	})

	// Лічильник звернень до GitHub API (допоможе стежити за лімітами)
	GitHubAPICalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "github_api_requests_total",
		Help: "Кількість запитів до GitHub API",
	}, []string{"status"}) // status може бути "success", "error", "rate_limited"
)
