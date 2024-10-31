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

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

FROM alpine:latest

RUN apk add --no-cache curl

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE="/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE=/app/configs/config.yaml
COPY configs/config.yaml ${CONFIG_FILE}

CMD ["sh", "-c", "until nc -z db 5432; do echo 'Waiting for database...'; sleep 2; done; /opt/calendar/calendar-app -config /app/configs/config.yaml"]
