.PHONY: build run clean test fmt vet lint docker-build docker-integration-test

VERSION = $(shell grep -oE 'version = "[^"]+"' .cz.toml | cut -d'"' -f2)

build:
	go build -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=$(VERSION)" -o bin/forte ./cmd/forte

run: build
	./bin/forte

clean:
	rm -f bin/forte

test:
	go test -timeout 30s -test.fullpath=true github.com/ericfortmeyer/forte/internal/version github.com/ericfortmeyer/forte/internal/help github.com/ericfortmeyer/forte/internal/fhs github.com/ericfortmeyer/forte/internal/deploy

fmt:
	go fmt github.com/ericfortmeyer/forte/internal/version github.com/ericfortmeyer/forte/internal/help github.com/ericfortmeyer/forte/internal/fhs github.com/ericfortmeyer/forte/internal/deploy

vet:
	go vet github.com/ericfortmeyer/forte/internal/version github.com/ericfortmeyer/forte/internal/help github.com/ericfortmeyer/forte/internal/fhs github.com/ericfortmeyer/forte/internal/deploy

lint: fmt vet
	@echo "Linting passed"

docker-build:
	docker build -qt forte:$(VERSION) .

docker-integration-test: build docker-build
	docker run --rm \
		-v $(shell pwd)/bin/forte:/usr/local/bin/forte \
		-v $(shell pwd):/usr/local/src \
		forte:$(VERSION) \
		bats -rpx ./integration_tests/
