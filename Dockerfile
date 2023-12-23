# syntax=docker/dockerfile:1
FROM golang:1.21-alpine as build
WORKDIR /app

# for Makefile
RUN apk add build-base
# for golang wire
RUN go install github.com/google/wire/cmd/wire@latest
# for golang swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN make ioc swagger
RUN go build -mod vendor -o ./.kanthor/kanthor -buildvcs=false

FROM alpine:3
WORKDIR /app

COPY --from=build /app/data ./data
COPY --from=build /app/configs.yaml ./configs.yaml

COPY --from=build /app/.kanthor/kanthor /usr/bin/kanthor

# debugging
EXPOSE 6060
# sdk
EXPOSE 8180
# portal
EXPOSE 8280