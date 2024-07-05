.PHONY: build
build:
	go build -v -o bin/orgonaut ./cmd/app

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: run
run:
	go mod tidy && go run cmd/app/main.go

.PHONY: docker_build
docker_build:
	docker build -t orgonaut:local -f local.Dockerfile .

.PHONY: docker_run
docker_run:
	docker run -d --name orgonaut orgonaut:local

.PHONY: docker_start
docker_start:
	docker start orgonaut

.PHONY: docker_stop
docker_stop:
	docker stop orgonaut

.PHONY: docker_remove
docker_remove:
	docker rm orgonaut

.PHONY: lint
lint:
	golangci-lint run


.DEFAULT_GOAL := build
