FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY listenerSvcApp /app/

CMD ["/app/listenerSvcApp"]
