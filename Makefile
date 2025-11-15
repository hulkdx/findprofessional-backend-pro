GO_VERSION := $(shell grep '^go ' professional-service/go.mod | awk '{print $$2}')

.PHONY: print-go-version
print-go-version:
	@echo "Go version from go.mod: $(GO_VERSION)"

.PHONY: deps
deps:
	@cd professional-service && \
	go mod download

.PHONY: test-deps
test-deps:
	@cd professional-service && \
	go get -t ./...

.PHONY: build
build: deps
	@cd professional-service && \
	go build -o ../build/app cmd/api/main.go

.PHONY: build-ttl
build-ttl: deps
	@cd professional-service && \
	go build -o ../build/app cmd/booking-holds-ttl/main.go

.PHONY: test
test: test-deps
	@cd professional-service && \
	go test ./...
	@git checkout professional-service/go.mod professional-service/go.sum

.PHONY: dev
dev:
	@cd dev-tools && \
	skaffold dev --port-forward 

.PHONY: clear-minikube-psql-cache
clear-minikube-psql-cache:
	PV=$$(kubectl get pvc data-postgresdb-0 -o jsonpath='{.spec.volumeName}'); \
	HOST=$$(kubectl get pv $$PV -o jsonpath='{.spec.hostPath.path}'); \
	minikube ssh -- "sudo rm -rf $$HOST"

.PHONY: docker-build
docker-build:
	docker build -f Dockerfile \
	--build-arg GO_VERSION=$(GO_VERSION) \
	--build-arg APP_CMD_PATH=$(APP_CMD_PATH) \
	-t $(DOCKER_IMAGE_NAME) .
