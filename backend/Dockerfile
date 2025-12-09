FROM golang:1.24.4-alpine AS builder
WORKDIR /src
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates netcat-openbsd
COPY --from=builder /app/server /server
EXPOSE 8085
ENTRYPOINT ["/bin/sh","-c","until nc -z postgres_db_redclass 5432 && nc -z redis_cache_redclass 6379; do echo waiting for postgres and redis; sleep 1; done; /server"]
