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
	"github-release-notifier/internal/grpcapi" // ДОДАНО: Твій gRPC пакет
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
	"github.com/soheilhy/cmux" // ДОДАНО: Мультиплексор
	"google.golang.org/grpc"   // ДОДАНО: Бібліотека gRPC
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

	// 1. Читаємо адресу Redis з оточення
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

	// 4. Читаємо GitHub токен з налаштувань Render
	githubToken := os.Getenv("GITHUB_TOKEN")

	// Передаємо прочитаний токен та клієнт Redis
	ghClient := github.NewClient(githubToken, rdb)

	handler := api.NewHandler(dbRepo, ghClient)

	emailNotifier := service.NewNotifier("smtp.gmail.com", "587", "tviy_email@gmail.com", "tviy_password")
	scanner := service.NewScanner(dbRepo, ghClient, emailNotifier, 5*time.Minute)

	scanner.Start()

	// --- Налаштування REST (Chi) ---
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

	// =====================================================================
	// МУЛЬТИПЛЕКСУВАННЯ (REST + gRPC НА ОДНОМУ ПОРТУ)
	// =====================================================================

	// 1. Створюємо основний слухач (Listener) на порту 8081
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Помилка створення слухача: %v", err)
	}

	// 2. Створюємо мультиплексор
	m := cmux.New(l)

	// 3. Правила сортування трафіку:
	// Якщо це HTTP/2 і в заголовках є application/grpc -> це gRPC
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// Усе інше (звичайні запити) відправляємо в наш HTTP/REST роутер
	httpL := m.Match(cmux.Any())

	// 4. Ініціалізуємо gRPC сервер
	grpcServer := grpc.NewServer()
	// Реєструємо наш сервіс. Передаємо dbRepo для роботи з базою
	grpcapi.RegisterNotifierServiceServer(grpcServer, grpcapi.NewGrpcHandler(dbRepo))

	// 5. Ініціалізуємо HTTP сервер (з нашим Chi роутером)
	httpServer := &http.Server{
		Handler: r,
	}

	// 6. Запускаємо gRPC сервер у фоні
	go func() {
		log.Println("gRPC-обробник готовий приймати запити")
		if err := grpcServer.Serve(grpcL); err != nil {
			log.Printf("gRPC сервер зупинився: %v", err)
		}
	}()

	// 7. Запускаємо HTTP сервер у фоні
	go func() {
		log.Println("REST-обробник готовий приймати запити")
		if err := httpServer.Serve(httpL); err != nil {
			log.Printf("HTTP сервер зупинився: %v", err)
		}
	}()

	// 8. Запускаємо сам мультиплексор (Він тримає програму відкритою)
	log.Println("🔥 Гібридний сервер (REST + gRPC) запущено на порту 8081")
	if err := m.Serve(); err != nil {
		log.Fatalf("Помилка мультиплексора: %v", err)
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
