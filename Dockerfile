# Етап 1: Збірка (Build)
FROM golang:1.25-alpine AS builder

# Встановлюємо робочу директорію
WORKDIR /app

# Копіюємо файли залежностей
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8081

CMD ["./main"]