package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github-release-notifier/internal/api"
	"github-release-notifier/internal/github"
	"github-release-notifier/internal/repository"
	"github-release-notifier/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Читаємо адресу бази з налаштувань Docker, або беремо локальну
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://postgres:secret@localhost:5432/notifier_db?sslmode=disable"
	}

	// Підключаємося до БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Помилка підключення до БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("БД не відповідає: %v", err)
	}

	// 2. Запуск міграцій
	runMigrations(db)

	// 3. Ініціалізація всіх компонентів системи
	dbRepo := repository.NewRepository(db)
	ghClient := github.NewClient("") // Поки без токена
	handler := api.NewHandler(dbRepo, ghClient)

	// Ініціалізація Notifier (пошта) та Scanner (фонова перевірка)
	emailNotifier := service.NewNotifier("smtp.gmail.com", "587", "tviy_email@gmail.com", "tviy_password")
	scanner := service.NewScanner(dbRepo, ghClient, emailNotifier, 5*time.Minute)

	// Запускаємо фоновий сканер
	scanner.Start()

	// 4. Налаштування роутера (HTTP-сервера)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// --- МАРШРУТИ ---

	// Публічний роут (без пароля), щоб перевіряти, чи живий сервер
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Захищена група роутів (тільки з паролем)
	r.Group(func(r chi.Router) {
		r.Use(api.APIKeyMiddleware) // Вмикаємо наш захист (Фейсконтроль)
		r.Post("/api/subscribe", handler.Subscribe)
	})
	// ----------------

	// 5. Запуск сервера на порту 8081
	log.Println("Сервер запущено на http://localhost:8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("Помилка запуску сервера: %v", err)
	}
}

// Функція міграцій
func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Помилка створення драйвера міграцій: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("Помилка ініціалізації міграцій: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Помилка виконання міграцій: %v", err)
	}

	log.Println("Міграції успішно виконано (таблиці створені)!")
}
