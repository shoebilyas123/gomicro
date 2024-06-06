FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY templates /app/templates
COPY mailApp /app/

CMD ["/app/mailApp"]
