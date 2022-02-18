FROM alpine:3.13.5

COPY ./bin/app /app/app

ENTRYPOINT ["/app/app"]
