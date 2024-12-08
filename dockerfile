# Start from the official Golang image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy the .env file into the container
# COPY .env .env

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the .env file from the builder stage to the root directory
# COPY --from=builder /app/.env .env

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]