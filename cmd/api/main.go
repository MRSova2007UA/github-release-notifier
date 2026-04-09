package main

import (
	"context"
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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func main() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://postgres:secret@localhost:5432/notifier_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Помилка ініціалізації БД: %v", err)
	}
	defer db.Close()

	// Намагаємося підключитися 5 разів з інтервалом у 2 секунди
	var pingErr error
	for i := 0; i < 5; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		log.Printf("БД ще не готова, чекаємо 2 секунди... (Спроба %d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		log.Fatalf("БД так і не відповіла після 5 спроб: %v", pingErr)
	}
	log.Println("Успішно підключено до PostgreSQL!")

	runMigrations(db)

	// 1. Читаємо адресу Redis з оточення (або ставимо локальну для тестів)
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// 2. Створюємо клієнт Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// 3. Перевіряємо, чи є зв'язок
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Помилка підключення до Redis: %v", err)
	}
	log.Println("Успішно підключено до Redis!")

	dbRepo := repository.NewRepository(db)

	// ВАЖЛИВО: Тепер ми передаємо rdb (клієнт Redis) у GitHub клієнт
	ghClient := github.NewClient("", rdb)
	handler := api.NewHandler(dbRepo, ghClient)

	emailNotifier := service.NewNotifier("smtp.gmail.com", "587", "tviy_email@gmail.com", "tviy_password")
	scanner := service.NewScanner(dbRepo, ghClient, emailNotifier, 5*time.Minute)

	scanner.Start()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	r.Handle("/metrics", promhttp.Handler())

	r.Group(func(r chi.Router) {
		r.Use(api.APIKeyMiddleware)
		r.Post("/api/subscribe", handler.Subscribe)
	})

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
