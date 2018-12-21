FROM alpine:3.8

ARG BINARY

ADD build/${BINARY} /usr/local/bin/nightwatch

ENTRYPOINT [ "nightwatch" ]