# Этап сборки
FROM golang:1.22 AS build

ENV BIN_FILE=/opt/calendar/calendar-app
ENV CODE_DIR=/go/src/app/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

# На выходе тонкий образ
FROM alpine:latest

# Установка curl для проверки сервиса
RUN apk add --no-cache curl

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE="/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

# Путь к конфигурационному файлу
ENV CONFIG_FILE=/app/configs/config.yaml
COPY configs/config.yaml ${CONFIG_FILE}

# Команда для проверки доступности базы данных перед запуском приложения
CMD ["sh", "-c", "until nc -z db 5432; do echo 'Waiting for database...'; sleep 2; done; /opt/calendar/calendar-app -config /app/configs/config.yaml"]
