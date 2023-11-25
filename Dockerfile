# syntax=docker/dockerfile:1
FROM kanthorlabs/base:latest as build
WORKDIR /app

COPY . .
RUN ls -la ./scripts/gen_ioc.sh
RUN make
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