package service

import (
	"log"
	"time"

	"github-release-notifier/internal/github"
	"github-release-notifier/internal/repository"
)

type Scanner struct {
	repo       *repository.Repository
	ghClient   *github.Client
	notifier   *Notifier
	pollPeriod time.Duration
}

func NewScanner(repo *repository.Repository, ghClient *github.Client, notifier *Notifier, period time.Duration) *Scanner {
	return &Scanner{
		repo:       repo,
		ghClient:   ghClient,
		notifier:   notifier,
		pollPeriod: period,
	}
}

// Start запускає сканер у нескінченному циклі
func (s *Scanner) Start() {
	ticker := time.NewTicker(s.pollPeriod)
	log.Printf("Сканер запущено. Інтервал перевірки: %v\n", s.pollPeriod)

	// Запускаємо горутину (фоновий процес)
	go func() {
		for {
			<-ticker.C // Чекаємо наступного тіку (наприклад, 5 хвилин)
			s.scan()
		}
	}()
}

func (s *Scanner) scan() {
	// 1. Отримуємо всі активні репозиторії
	repos, err := s.repo.GetActiveRepositories()
	if err != nil {
		log.Printf("Помилка отримання репозиторіїв для сканування: %v", err)
		return
	}

	for _, repo := range repos {
		repoID := repo["id"]
		repoName := repo["name"]
		lastSeenTag := repo["last_seen_tag"]

		// 2. Запитуємо GitHub про останній реліз
		latestTag, err := s.ghClient.GetLatestRelease(repoName)
		if err != nil {
			log.Printf("Помилка перевірки релізу для %s: %v", repoName, err)
			continue
		}

		// 3. Порівнюємо теги. Якщо новий — діємо!
		if latestTag != "" && latestTag != lastSeenTag {
			log.Printf("Знайдено новий реліз для %s: %s (було %s)", repoName, latestTag, lastSeenTag)

			// Отримуємо email-и всіх, хто підписався на цей репозиторій
			emails, err := s.repo.GetSubscribersForRepo(repoID)
			if err != nil {
				log.Printf("Помилка отримання підписників для %s: %v", repoName, err)
				continue
			}

			// Відправляємо листи
			if err := s.notifier.SendReleaseEmail(emails, repoName, latestTag); err != nil {
				log.Printf("Помилка відправки листів для %s: %v", repoName, err)
				continue
			}

			// Оновлюємо тег у базі даних, щоб не відправити лист двічі
			if err := s.repo.UpdateLastSeenTag(repoID, latestTag); err != nil {
				log.Printf("Помилка оновлення тегу в БД для %s: %v", repoName, err)
			}
		}
	}
}
