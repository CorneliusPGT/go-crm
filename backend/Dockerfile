FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o crm ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/crm .
COPY --from=builder /app/internal/migrations ./internal/migrations

ENV POSTGRES_DSN="postgres://postgres:postgres@db:5432/crm?sslmode=disable"

EXPOSE 8080

CMD ["./crm"]