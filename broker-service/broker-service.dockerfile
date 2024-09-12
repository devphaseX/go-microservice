FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# Start a new stage from scratch
FROM alpine:latest

RUN mkdir /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/brokerApp /app

EXPOSE 5001

WORKDIR /app

CMD ["./brokerApp"]
