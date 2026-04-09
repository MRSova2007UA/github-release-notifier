# Етап 1: Збірка (Build)
FROM golang:1.25-alpine AS builder

# Встановлюємо робочу директорію
WORKDIR /app

# Копіюємо файли залежностей
COPY go.mod go.sum ./
RUN go mod download

# Копіюємо весь код
COPY . .

# Компілюємо бінарний файл
RUN go build -o main ./cmd/api/main.go

# Етап 2: Запуск (Run)
FROM alpine:latest

WORKDIR /app

# Копіюємо тільки скомпільований файл та папку з міграціями
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Відкриваємо порт
EXPOSE 8081

# Запускаємо програму
CMD ["./main"]