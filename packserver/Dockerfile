# Use an official Go image as the build stage.
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code.
COPY . .

# Run tests. If any test fails, the build will fail.
RUN go test -v ./...

# Build the binary for Linux.
RUN CGO_ENABLED=0 GOOS=linux go build -o packserver .

# Use a minimal base image for the final container.
FROM alpine:latest

# Set working directory in final container.
WORKDIR /root/

# Copy the built binary from the builder stage.
COPY --from=builder /app/packserver .

# Expose port 5000.
EXPOSE 5000

# Command to run the binary.
CMD ["./packserver"]
