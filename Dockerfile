FROM alpine:3.13.5

COPY ./bin/app /app/app

EXPOSE 7000

ENTRYPOINT ["/app/app"]
