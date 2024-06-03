# build the go app
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# Run the tiny app in a docker container
FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/brokerApp /app/

RUN ["/app/brokerApp"]