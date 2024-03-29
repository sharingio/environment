FROM alpine:3.14 AS extras
RUN apk add --no-cache tzdata ca-certificates
RUN adduser -D user

FROM scratch AS base
WORKDIR /app
ENV PATH=/app/bin
COPY --from=extras /etc/passwd /etc/passwd
COPY --from=extras /etc/group /etc/group
COPY --from=extras /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=extras /etc/ssl /etc/ssl
USER user

FROM golang:1.17.0-alpine3.14 AS build
WORKDIR /app
COPY pkg /app/pkg
COPY cmd /app/cmd
COPY go.* *.go /app/
ARG GOARCH=""
RUN CGO_ENABLED=0 GOOS=linux GOARCH="$GOARCH" go build \
  -a \
  -installsuffix cgo \
  -ldflags "-extldflags '-static' -s -w" \
  -o bin/environment-exporter \
  cmd/environment-exporter/main.go

FROM base AS final
COPY --from=build /app/bin/environment-exporter /app/bin/environment-exporter
EXPOSE 10093
ENTRYPOINT ["/app/bin/environment-exporter"]
