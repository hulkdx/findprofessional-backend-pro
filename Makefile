.PHONY: dev
dev:
	cd local-development && \
	skaffold dev --port-forward 

.PHONY: clear-minikube-psql-cache
clear-minikube-psql-cache:
	eval $$(minikube docker-env); \
	docker volume rm --force psql_cache
