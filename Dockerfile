FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build -o feed ./cmd/feed
RUN CGO_ENABLED=1 go build -o migrator ./cmd/migrator

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/feed ./feed
COPY --from=builder /app/migrator ./migrator
COPY ./config/dev.yaml ./config/dev.yaml

ENTRYPOINT ["./feed", "--config=./config/dev.yaml"]