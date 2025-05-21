FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY services/products /app/services/products
COPY services/auth/pkg/auth /app/services/auth/pkg/auth
COPY clients /app/clients
COPY protos /app/protos
COPY temporal-app /app/temporal-app

WORKDIR /app/services/products

RUN go mod download

WORKDIR /app/services/products
RUN go build -o app ./cmd/cron/main.go


FROM golang:1.23-alpine as release

WORKDIR /app

COPY --from=builder /app/services/products/app .
COPY --from=builder /app/services/products/config.yml .
COPY --from=builder /app/services/products/migrations /app/migrations

CMD ["./app"]