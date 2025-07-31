# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY resource ./resource
EXPOSE 8080
CMD ["./app"]
