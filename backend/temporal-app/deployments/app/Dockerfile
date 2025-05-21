FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY temporal-app /app/temporal-app
COPY clients /app/clients
COPY protos /app/protos

WORKDIR /app/temporal-app

RUN go mod download

WORKDIR /app/temporal-app
RUN go build -o app ./cmd/app/main.go


FROM golang:1.23-alpine as release

WORKDIR /app

COPY --from=builder /app/temporal-app/app .
COPY --from=builder /app/temporal-app/config.yml .

CMD ["./app"]