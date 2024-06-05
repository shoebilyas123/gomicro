FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY logsvcApp /app/

CMD ["/app/logsvcApp"]
