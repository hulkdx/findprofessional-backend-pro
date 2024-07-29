FROM golang:1.22.5-alpine AS builder
WORKDIR /src
COPY professional-service .
RUN go mod download
RUN go build \
    -ldflags "-w -s"\
    -buildvcs=false \
    -o app \
    cmd/api/main.go

# -----------------------------------------------------------------------------
FROM alpine:latest
COPY --from=builder /src/app app
ENTRYPOINT [ "./app" ]
