.PHONY: docs
docs:
	@swag init -g ./cmd/main.go && swag fmt

.PHONY: build
build:
	@go build -o ./bin/main ./cmd/

.PHONY: dev
dev:
	@go run ./cmd/

.PHONY: run
run:
	@./bin/main

.PHONY: build-docker
build-docker:
	@docker build -t autodeploy-api:latest .