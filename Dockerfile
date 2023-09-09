# syntax=docker/dockerfile:1
FROM scrapnode/kanthor-base:latest as build
WORKDIR /app

COPY . .
RUN make
RUN go build -mod vendor -o ./.kanthor/kanthor -buildvcs=false

FROM alpine:3
WORKDIR /app

COPY --from=build /app/data ./data
COPY --from=build /app/migration ./migration
COPY --from=build /app/configs.yaml ./configs.yaml

COPY --from=build /app/.kanthor/kanthor /usr/bin/kanthor

# sdkapi
EXPOSE 8180
# portalapi
EXPOSE 8280