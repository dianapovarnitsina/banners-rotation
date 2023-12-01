# Собираем в гошке
FROM golang:1.20.5 as build

ENV BIN_FILE /opt/banner/banner-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY build/banner ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
COPY . /go/src/
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/banner/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="banner"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/banner/banner-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/banner/banner_config.yaml
COPY ./configs/banner_config.yaml ${CONFIG_FILE}

ENV CONFIG_MIGRATION /etc/migrations
COPY ./migrations ${CONFIG_MIGRATION}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
