#!/usr/bin/make

.DEFAULT_GOAL := help

APP_TAG := 0.1.2


help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo "\n  Allowed for overriding next properties:\n\n\
		Usage example:\n\
		make build"

build-ghcr:
	docker build --rm -f Dockerfile -t ghcr.io/aak74/task-exporter:latest -t ghcr.io/aak74/task-exporter:$(APP_TAG) .

push-ghcr:
	docker push ghcr.io/aak74/task-exporter:latest
	docker push ghcr.io/aak74/task-exporter:$(APP_TAG)
