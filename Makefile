API_BIN := "./bin/banner"
DOCKER_IMG="banner:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(API_BIN) -ldflags "$(LDFLAGS)" ./cmd/banner


run: build
	$(API_BIN) -config ./configs/banner_config.yaml


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

generate:
	rm -rf internal/server/pb
	mkdir -p internal/server/pb

	protoc \
        --proto_path=api/ \
        --go_out=internal/server/pb \
        --go-grpc_out=internal/server/pb \
        api/*.proto

test:
	go test -race ./internal/...

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/banner/Dockerfile .




