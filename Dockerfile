# syntax=docker/dockerfile:1
FROM golang:1.20-alpine as build
WORKDIR /app

RUN apk add build-base
RUN go install github.com/google/wire/cmd/wire@latest

COPY . .
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download

RUN make gen-go
RUN go build -mod vendor -o ./.kanthor/kanthor -buildvcs=false

FROM alpine:3
WORKDIR /app

COPY --from=build /app/data ./data
COPY --from=build /app/migration ./migration
COPY --from=build /app/configs.yaml ./configs.yaml

COPY --from=build /app/.kanthor/kanthor ./kanthor
COPY --from=build /app/docker/entrypoint.sh ./entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# portalapi
EXPOSE 8080,9090
# sdkapi
EXPOSE 8180,9190
# scheduler
EXPOSE 8280,9290
# dispatcher
EXPOSE 8380,9390

ENTRYPOINT ["/app/entrypoint.sh"]