# syntax=docker/dockerfile:1
FROM golang:1.20-alpine as build
WORKDIR /app

COPY . .

# for Makefile
RUN apk add build-base
# for golang wire
RUN go install github.com/google/wire/cmd/wire@latest
# for golang swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN make
RUN go build -mod vendor -o ./.kanthor/kanthor -buildvcs=false

FROM alpine:3
WORKDIR /app

COPY --from=build /app/data ./data
COPY --from=build /app/migration ./migration
COPY --from=build /app/configs.yaml ./configs.yaml

COPY --from=build /app/.kanthor/kanthor ./kanthor
COPY --from=build /app/docker/entrypoint.sh ./entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# sdkapi
EXPOSE 8080
# portalapi
EXPOSE 8180
# scheduler
EXPOSE 8280
# dispatcher
EXPOSE 8380

ENTRYPOINT ["/app/entrypoint.sh"]