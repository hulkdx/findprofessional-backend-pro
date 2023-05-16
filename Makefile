.PHONY: deps
deps:
	@cd professional-service && \
	go mod download

.PHONY: build
build: deps
	@cd professional-service && \
	go build -o ../build/app cmd/api/main.go

.PHONY: run
run: build
	@./build/app

.PHONY: test
test:
	@cd professional-service && \
	go test ./...

.PHONY: dev
dev:
	@cd dev-tools && \
	skaffold dev --port-forward 

.PHONY: clear-minikube-psql-cache
clear-minikube-psql-cache:
	@eval $$(minikube docker-env); \
	docker volume rm --force psql_cache
