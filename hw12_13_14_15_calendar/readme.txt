docker собирается из корневой папки командой docker buildx build -f deployments/build/Dockerfile .

docker exec -it calendar_service psql -U user -h db -d calendar       входит
docker exec -it calendar_service psql -h postgres_db -U user -d calendar
docker exec -it postgres_db psql -U user -d calendar


CMD ["/bin/sh", "-c", "sleep 15 && /opt/calendar/calendar-app -config /app/configs/config.yaml"]
CMD ["${BIN_FILE}", "-config", "${CONFIG_FILE}"]

docker-compose down
docker network prune  # Внимание: это удалит все неиспользуемые сети
docker-compose up -d

docker network ls

docker network inspect calendar_network

docker exec -it calendar_service ping db

docker ps            проверьте состояние всех контейнеров после запуска:

docker logs postgres_db              можете просмотреть логи контейнера postgres_db




FROM golang:1.22 AS build

ENV BIN_FILE=/opt/calendar/calendar-app
ENV CODE_DIR=/go/src/

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
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE="/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

# Измените этот путь на правильный
ENV CONFIG_FILE=/app/configs/config.yaml
COPY ../../configs/config.yaml ${CONFIG_FILE}

CMD ["/bin/sh", "-c", "/opt/calendar/calendar-app -config /app/configs/config.yaml"]

