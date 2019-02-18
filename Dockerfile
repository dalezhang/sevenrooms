FROM alpine:3.6
MAINTAINER Dale dale.zhang@bindo.com
RUN apk update && \
    apk add ca-certificates musl-dev && \
    apk add tzdata
COPY main /app/main
COPY config /app/config

CMD cd /app && ./main
