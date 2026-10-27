# Stage 1: Build binary
FROM golang:alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0

RUN apk add --no-cache git ca-certificates tzdata

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -ldflags="-w -s" -o /app/musashi-api ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/musashi-api /app/musashi-api

EXPOSE 8080

CMD ["/app/musashi-api"]

