![Build](https://img.shields.io/github/actions/workflow/status/hulkdx/findprofessional-backend-pro/push.yml?branch=main)
[![Docker Status](https://badgen.net/docker/size/hulkdx/ff-pro/v1/amd64?icon=docker&label=docker&url)](https://hub.docker.com/repository/docker/hulkdx/ff-pro)
[![Golang](https://img.shields.io/badge/golang-1.22.5-blue.svg?logo=go)](https://go.dev/)

# Professional microservice

## Development
Requirements
- [Go](https://go.dev)
- [Docker](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) (or any other localhost kubernetes)
- [Skaffold](https://skaffold.dev/docs/install/)

To start development, run `make dev`

## How to upgrade go version
Check [this](https://github.com/hulkdx/findprofessional-backend-pro/commit/af66c3d722d3553ff01137072d7c5077471415a7) commit

## TODO
- Host swagger docs to some url
