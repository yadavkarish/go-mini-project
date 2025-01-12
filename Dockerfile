# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code to the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/app/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Add CA certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your app listens on
EXPOSE 8081

# Run the application
ENTRYPOINT ["./main"]
