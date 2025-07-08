FROM golang:1.23.4-alpine3.19 AS builder

RUN apk --update add git make openssh

RUN addgroup -g 1001 -S builder && \
    adduser -u 1001 -S builder -G builder
USER builder:builder

ENV APP_HOME /home/builder
WORKDIR $APP_HOME

COPY --chown=builder:builder . $APP_HOME/

# build executable to ./bin/app
RUN make build-native

FROM alpine:3.19

RUN addgroup -g 1001 -S runner && \
    adduser -u 1001 -S runner -G runner

ENV BUILD_DIR /home/builder
ENV APP_HOME /home/runner/bin
# create dummy .env file
RUN mkdir -p $APP_HOME && touch $APP_HOME/.env

USER runner:runner
WORKDIR $APP_HOME

COPY --chown=runner:runner --from=builder $BUILD_DIR/bin .
COPY --chown=runner:runner --from=builder $BUILD_DIR/resources ./resources

ENTRYPOINT ["./app"]

# override CMD with either "http" or "kafka -t <topic-name>"
CMD ["--help"]
