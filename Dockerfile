FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY . .

CMD ["/app/bin/server"]