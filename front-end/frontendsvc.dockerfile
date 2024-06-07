FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY frontendApp /app/

CMD ["/app/frontendApp"]