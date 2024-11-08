FROM golang:1.22 AS build

RUN apt-get update && apt-get install -y \
     gcc \
     libc-dev \
     git \
     && rm -rf /var/lib/apt/lists/*

ENV BIN_FILE=/opt/calendar/calendar-app
ENV SCHEDULER_BIN_FILE=/opt/calendar/calendar-scheduler-app
ENV SENDER_BIN_FILE=/opt/calendar/calendar-sender-app
ENV CODE_DIR=/go/src/app/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${SCHEDULER_BIN_FILE} cmd/calendar_scheduler/*

RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${SENDER_BIN_FILE} cmd/calendar_sender/*

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

FROM debian:latest

RUN apt-get update && apt-get install -y \
    curl \
    postgresql-client \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

COPY --from=build /opt/calendar/calendar-app /opt/calendar/calendar-app
COPY --from=build /opt/calendar/calendar-scheduler-app /opt/calendar/calendar-scheduler-app
COPY --from=build /opt/calendar/calendar-sender-app /opt/calendar/calendar-sender-app
COPY --from=build /go/bin/grpcurl /usr/local/bin/grpcurl

ENV CONFIG_FILE=/app/configs/calendar_config.yaml
COPY configs/calendar_config.yaml ${CONFIG_FILE}
ENV SCHEDULER_CONFIG_FILE=/app/configs/scheduler_config.yaml
COPY configs/scheduler_config.yaml ${SCHEDULER_CONFIG_FILE}
ENV SENDER_CONFIG_FILE=/app/configs/sender_config.yaml
COPY configs/sender_config.yaml ${SENDER_CONFIG_FILE}

COPY ./migrations /migrations

CMD ["sh", "-c", "until nc -z db 5432; do echo 'Waiting for database...'; sleep 2; done; /opt/calendar/calendar-app -config /app/configs/calendar_config.yaml"]