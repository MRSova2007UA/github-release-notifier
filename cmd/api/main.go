package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github-release-notifier/internal/api"
	"github-release-notifier/internal/github"
	"github-release-notifier/internal/grpcapi"
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
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

func main() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://postgres:secret@localhost:5432/notifier_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Помилка ініціалізації бази даних: %v", err)
	}
	defer db.Close()

	var pingErr error
	for i := 0; i < 5; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		log.Printf("Очікування готовності БД... (Спроба %d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		log.Fatalf("Не вдалося підключитися до БД після 5 спроб: %v", pingErr)
	}
	log.Println("Успішне підключення до PostgreSQL")

	runMigrations(db)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Помилка підключення до Redis: %v", err)
	}
	log.Println("Успішне підключення до Redis")

	dbRepo := repository.NewRepository(db)
	githubToken := os.Getenv("GITHUB_TOKEN")
	ghClient := github.NewClient(githubToken, rdb)

	handler := api.NewHandler(dbRepo, ghClient)

	// Налаштування системи сповіщень через SMTP
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	emailNotifier := service.NewNotifier("smtp.gmail.com", "587", smtpUser, smtpPass)

	scanner := service.NewScanner(dbRepo, ghClient, emailNotifier, 5*time.Minute)
	scanner.Start()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Помилка прослуховування порту 8080: %v", err)
	}

	m := cmux.New(l)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	grpcServer := grpc.NewServer()
	grpcapi.RegisterNotifierServiceServer(grpcServer, grpcapi.NewGrpcHandler(dbRepo))

	httpServer := &http.Server{
		Handler: r,
	}

	go func() {
		log.Println("gRPC обробник готовий до роботи")
		if err := grpcServer.Serve(grpcL); err != nil {
			log.Printf("Помилка gRPC сервера: %v", err)
		}
	}()

	go func() {
		log.Println("REST обробник готовий до роботи")
		if err := httpServer.Serve(httpL); err != nil {
			log.Printf("Помилка HTTP сервера: %v", err)
		}
	}()

	log.Println("Гібридний сервер (REST + gRPC) запущено на порту 8080")
	if err := m.Serve(); err != nil {
		log.Fatalf("Помилка мультиплексора: %v", err)
	}
}

// runMigrations виконує оновлення схеми бази даних за допомогою міграційних файлів
func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Помилка драйвера міграцій: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("Помилка ініціалізації міграцій: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Помилка виконання міграцій: %v", err)
	}

	log.Println("Міграції успішно застосовані")
}
