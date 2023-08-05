# syntax=docker/dockerfile:1
FROM golang:1.20-alpine as build
WORKDIR /app

RUN apk add build-base
RUN go install github.com/google/wire/cmd/wire@latest

COPY . .

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

ENTRYPOINT ["/app/entrypoint.sh"]