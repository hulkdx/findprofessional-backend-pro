![Build](https://img.shields.io/github/actions/workflow/status/hulkdx/findprofessional-backend-pro/push.yaml?branch=main)
[![Docker Status](https://badgen.net/docker/size/hulkdx/ff-pro/v1/amd64?icon=docker&label=docker&url)](https://hub.docker.com/repository/docker/hulkdx/ff-pro)
[![TTL Docker Status](https://badgen.net/docker/size/hulkdx/ff-pro-ttl/v1/amd64?icon=docker&label=docker%20ttl&url)](https://hub.docker.com/repository/docker/hulkdx/ff-pro-ttl)
[![Golang](https://img.shields.io/badge/golang-1.23.0-blue.svg?logo=go)](https://go.dev/)

# Professional microservice

## Development
Requirements
- [Go](https://go.dev)
- [Docker](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) (or any other localhost kubernetes)
- [Skaffold](https://skaffold.dev/docs/install/)

### How to start local development
- Start docker
- Start Minikube
```sh
minikube start
```
- start skaffold
```sh
make dev
```

To start development, run `make dev`

## How to upgrade go version
Check [this](https://github.com/hulkdx/findprofessional-backend-pro/commit/af66c3d722d3553ff01137072d7c5077471415a7) commit

## TODO
- Host swagger docs to some url

# Booking Holds TTL
Deletes expired rows from booking_holds (and their booking_hold_items via ON DELETE CASCADE) on a schedule so that previously held slots become visible/available again.

This tool runs as a short-lived job with Kubernetes CronJob. It performs one sweep and exits with a success/failure code.

## Why it exists
- GET /professional endpoints hide a slot if there’s a non-expired hold.
- When a client abandons payment or a flow fails, the hold expires.
- To prevent an expired hold from blocking a slot, we periodically purge expired holds.
- Purging only the parent booking_holds ensures items disappear automatically.
