FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY services/orders /app/services/orders
COPY clients /app/clients
COPY protos /app/protos

WORKDIR /app/services/orders

RUN go mod download

WORKDIR /app/services/orders
RUN go build -o app ./cmd/app/main.go


FROM golang:1.23-alpine as release

WORKDIR /app

COPY --from=builder /app/services/orders/app .
COPY --from=builder /app/services/orders/config.yml .
COPY --from=builder /app/services/orders/migrations /app/migrations

CMD ["./app"]