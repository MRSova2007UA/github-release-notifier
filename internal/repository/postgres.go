package repository

import (
	"database/sql"
	"fmt"
)

// Repository — структура для роботи з БД
type Repository struct {
	db *sql.DB
}

// NewRepository створює новий екземпляр
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// SubscribeUser зберігає email та репозиторій, а потім зв'язує їх
func (r *Repository) SubscribeUser(email, repoName, latestTag string) error {
	// Починаємо транзакцію (щоб або все збереглося, або нічого)
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Якщо щось піде не так, транзакція відкотиться
	defer tx.Rollback()

	// 1. Додаємо користувача (або отримуємо його ID, якщо він вже є)
	var subscriberID int
	err = tx.QueryRow(`
		INSERT INTO subscribers (email) 
		VALUES ($1) 
		ON CONFLICT (email) DO UPDATE SET email=EXCLUDED.email 
		RETURNING id`, email).Scan(&subscriberID)
	if err != nil {
		return fmt.Errorf("помилка збереження підписника: %v", err)
	}

	// 2. Додаємо репозиторій (або отримуємо його ID, якщо він вже є)
	var repoID int
	err = tx.QueryRow(`
		INSERT INTO repositories (name, last_seen_tag) 
		VALUES ($1, $2) 
		ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name 
		RETURNING id`, repoName, latestTag).Scan(&repoID)
	if err != nil {
		return fmt.Errorf("помилка збереження репозиторію: %v", err)
	}

	// 3. Зв'язуємо їх у таблиці subscriptions (ігноруємо помилку, якщо підписка вже існує)
	_, err = tx.Exec(`
		INSERT INTO subscriptions (subscriber_id, repository_id) 
		VALUES ($1, $2) 
		ON CONFLICT DO NOTHING`, subscriberID, repoID)
	if err != nil {
		return fmt.Errorf("помилка збереження підписки: %v", err)
	}

	// Підтверджуємо транзакцію
	return tx.Commit()
}

// GetActiveRepositories повертає всі репозиторії, на які хтось підписаний
func (r *Repository) GetActiveRepositories() ([]map[string]string, error) {
	rows, err := r.db.Query(`
		SELECT id, name, last_seen_tag 
		FROM repositories 
		WHERE id IN (SELECT DISTINCT repository_id FROM subscriptions)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []map[string]string
	for rows.Next() {
		var id, name, tag string
		if err := rows.Scan(&id, &name, &tag); err != nil {
			continue
		}
		repos = append(repos, map[string]string{
			"id":            id,
			"name":          name,
			"last_seen_tag": tag,
		})
	}
	return repos, nil
}

// UpdateLastSeenTag оновлює тег після знаходження нового релізу
func (r *Repository) UpdateLastSeenTag(repoID, newTag string) error {
	_, err := r.db.Exec(`UPDATE repositories SET last_seen_tag = $1 WHERE id = $2`, newTag, repoID)
	return err
}

// GetSubscribersForRepo повертає список email-ів для конкретного репозиторію
func (r *Repository) GetSubscribersForRepo(repoID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT s.email 
		FROM subscribers s
		JOIN subscriptions sub ON s.id = sub.subscriber_id
		WHERE sub.repository_id = $1
	`, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			continue
		}
		emails = append(emails, email)
	}
	return emails, nil
}
