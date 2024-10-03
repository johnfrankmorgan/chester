FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./*.go ./
COPY ./embed ./embed

RUN go build -o /build/chester .

FROM lichessbotdevs/lichess-bot:alpine

COPY --from=builder /build/chester /lichess-bot/engines/chester

COPY ./docker/entrypoint.sh /entrypoint.sh
COPY ./docker/config.yml /lichess-bot/config.yml

ENTRYPOINT ["/entrypoint.sh"]
