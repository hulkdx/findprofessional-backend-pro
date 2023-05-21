FROM golang:1.20.4-alpine AS builder
WORKDIR /src
COPY professional-service .
RUN go mod download
RUN go build \
    -ldflags "-w -s"\
    -buildvcs=false \
    -o app \
    cmd/api/main.go
# -----------------------------------------------------------------------------
FROM alpine:3.16
COPY --from=builder /src/app app
ENTRYPOINT [ "./app" ]
