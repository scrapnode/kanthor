# syntax=docker/dockerfile:1
FROM namely/protoc-all:1.51_1 as protobuild
WORKDIR /app

COPY . .
RUN apt-get -y install make
RUN make gen-proto

FROM golang:1.20-alpine as build
WORKDIR /app

RUN apk add build-base
RUN go install github.com/google/wire/cmd/wire@latest

COPY --from=protobuild /app .
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download

RUN make gen-go
RUN go build -o ./.kanthor/kanthor -buildvcs=false

FROM alpine:3
WORKDIR /app

COPY --from=build /app/data ./data
COPY --from=build /app/migration ./migration
COPY --from=build /app/configs.yaml ./configs.yaml

COPY --from=build /app/docker/entrypoint.sh ./entrypoint.sh
RUN chmod +x /app/entrypoint.sh

COPY --from=build /app/.kanthor/kanthor ./kanthor

# controlplane
EXPOSE 8080,8081,9090
# dataplane
EXPOSE 8180,8181,9190
# scheduler
EXPOSE 8280,8281,9290
# dispatcher
EXPOSE 8380,8381,9390

ENTRYPOINT ["/app/entrypoint.sh"]