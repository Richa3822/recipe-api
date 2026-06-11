# Stage 1 — Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN go build -o main ./cmd/main.go

# Stage 2 — Run with a tiny image
FROM alpine:latest

WORKDIR /app

# Copy only the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]