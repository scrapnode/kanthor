# syntax=docker/dockerfile:1
FROM golang:1.21-alpine
WORKDIR /app

# for Makefile
RUN apk add build-base
# for golang wire
RUN go install github.com/google/wire/cmd/wire@latest
# for golang swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest
