FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY services/auth /app/services/auth
COPY clients /app/clients
COPY protos /app/protos
COPY temporal-app /app/temporal-app

WORKDIR /app/services/auth

RUN go mod download

WORKDIR /app/services/auth
RUN go build -o app ./cmd/app/main.go


FROM golang:1.23-alpine as release

WORKDIR /app

COPY --from=builder /app/services/auth/app .
COPY --from=builder /app/services/auth/config.yml .
COPY --from=builder /app/services/auth/migrations /app/migrations

CMD ["./app"]