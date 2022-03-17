FROM alpine:latest

COPY ./bin/app /app/app

ENTRYPOINT ["/app/app"]
