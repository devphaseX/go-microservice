FROM golang:1.23.1-alpine AS builder

RUN mkdir /app

COPY . ./app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp


FROM alphine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

EXPOSE 5001

WORKDIR /app

CMD ["./app/brokerApp"]
