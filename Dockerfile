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

COPY --from=build /app/migration ./migration
COPY --from=build /app/configs.yaml ./configs.yaml

COPY --from=build /app/docker/entrypoint.sh ./entrypoint.sh
RUN chmod +x /app/entrypoint.sh

COPY --from=build /app/.kanthor/kanthor ./kanthor

# HTTP for setup services
EXPOSE 8080-8089
# HTTP for dataplane service
EXPOSE 8180-8189
# HTTP for scheduler service
EXPOSE 8280-8289
# HTTP for dispatcher service
EXPOSE 8380-8389

# gRPC for setup services
EXPOSE 9080-9089
# gRPC for dataplane service
EXPOSE 9180-9189
# gRPC for scheduler service
EXPOSE 9280-9289
# gRPC for dispatcher service
EXPOSE 9380-9389

ENTRYPOINT ["/app/entrypoint.sh"]