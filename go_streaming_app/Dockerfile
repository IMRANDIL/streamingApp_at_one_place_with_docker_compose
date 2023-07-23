# Stage 1: Build the Go application
FROM golang:1.17-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code to the container
COPY . .

# Set the working directory to the cmd/api directory
WORKDIR /app/cmd/api

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Create the final Docker image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary built in the first stage
COPY --from=builder /app/cmd/api/app .

# Expose the port the application is listening on
EXPOSE 8080

# Run the Go application
CMD ["./app"]
