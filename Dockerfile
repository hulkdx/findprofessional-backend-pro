# --------------------------------  Build  -----------------------------------------
FROM golang:1.23.0-alpine AS builder
WORKDIR /src

COPY professional-service/go.mod professional-service/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY professional-service ./

ENV CGO_ENABLED=0 GOOS=linux
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -trimpath \
    -ldflags "-s -w"\
    -buildvcs=false \
    -o /out/app ./cmd/api

# --------------------------------  Runtime  -----------------------------------------
FROM alpine:latest

RUN adduser -D -H -u 10001 appuser
WORKDIR /home/appuser

COPY --from=builder /out/app app
USER appuser
ENTRYPOINT [ "./app" ]
