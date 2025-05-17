# Stage 1: Build the Go app
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy Go mod files and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o main .

# Stage 2: Minimal final image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

# Create a mountable output directory
VOLUME ["/app/outputs"]

EXPOSE 8080
CMD ["./main"]

