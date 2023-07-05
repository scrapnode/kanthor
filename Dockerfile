# syntax=docker/dockerfile:1
FROM golang:1.20-alpine as build
WORKDIR /app

RUN apk add build-base
COPY . .
RUN go mod download
RUN go build -o ./.kanthor/kanthor -buildvcs=false

FROM alpine:3
WORKDIR /app

COPY --from=build /app/migration ./migration
COPY --from=build /app/configs.yaml ./configs.yaml

COPY --from=build /app/docker/entrypoint.sh ./entrypoint.sh
RUN chmod +x /app/entrypoint.sh

COPY --from=build /app/.kanthor/kanthor ./kanthor

EXPOSE 8080-30010
ENTRYPOINT ["/app/entrypoint.sh"]