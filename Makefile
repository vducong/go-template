.PHONY: install setup-hooks build run test lint gen gen-proto docker-build docker-run

APP_NAME = gotmpl
DOCKER_IMAGE = $(APP_NAME):latest

install:
	@echo "Installing dependencies..."
	go mod tidy
	brew install buf golangci-lint
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@$(MAKE) setup-hooks

setup-hooks:
	chmod +x scripts/setup-hooks.sh
	./scripts/setup-hooks.sh

build:
	go build -o bin/app ./cmd/app

build-cli:
	go build -o bin/cli ./cmd/cli

run:
	go run cmd/app/main.go

test:
	@echo "Running tests..."
	go test -v ./...

lint:
	@echo "Running linters..."
	golangci-lint run

gen:
	go generate ./...

gen-proto:
	@echo "Generating Protobuf files..."
	buf generate

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)