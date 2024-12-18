FROM golang:1.21.1-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o app .
FROM alpine:latest
RUN apk add --no-cache libpq
COPY --from=builder /app/app /app/
COPY .env .env
WORKDIR /app
CMD ["./app"]
