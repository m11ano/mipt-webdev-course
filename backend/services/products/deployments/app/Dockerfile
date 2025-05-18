FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY services/products /app/services/products
COPY clients /app/clients

WORKDIR /app/services/products

RUN go mod download

WORKDIR /app/services/products
RUN go build -o app ./cmd/app/main.go


FROM golang:1.23-alpine as release

WORKDIR /app

COPY --from=builder /app/services/products/app .
COPY --from=builder /app/services/products/config.yml .
COPY --from=builder /app/services/products/migrations /app/migrations

CMD ["./app"]