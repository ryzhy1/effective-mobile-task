# Этап сборки
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go get -u github.com/swaggo/swag

# Копирование исходного кода
COPY . .

# Генерация Swagger документации
RUN swag init -g cmd/main.go

# Сборка приложения
RUN go build -o /app/bin/app ./cmd/main.go

# Этап запуска
FROM alpine:3.18 AS runner

WORKDIR /app

# Копирование бинарника приложения
COPY --from=builder /app/bin/app /app/app

# Копирование миграций
COPY migrations /app/migrations

# Копирование утилиты goose
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копирование Swagger документации
COPY --from=builder /app/docs /app/docs

# Копирование конфигурационного файла
COPY .env.local /app/.env.local

# Команда запуска
CMD ["sh", "-c", "goose -dir /app/migrations postgres 'user=postgres password=postgres dbname=postgres host=db port=5432 sslmode=disable' up && ./app"]
