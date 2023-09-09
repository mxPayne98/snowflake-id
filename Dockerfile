# Start from a base Go image
FROM golang:1.20 AS builder

# Set the working directory in Docker
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/main .

# Command to run the application
CMD ["./main"]
