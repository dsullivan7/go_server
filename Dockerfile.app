FROM alpine:3.13.5

RUN \
  apk update && \
  apk upgrade && \
  apk add chromium

COPY ./bin/app /app/app

ENTRYPOINT ["/app/app"]
