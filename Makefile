.PHONY: build run clean test fmt vet lint docker-build docker-integration-test example-php-deployment

VERSION = $(shell grep -oE 'version = "[^"]+"' .cz.toml | cut -d'"' -f2 || echo "0.0.0")

build:
	go build -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=$(VERSION)" -o bin/forte ./cmd/forte

run: build
	./bin/forte

clean:
	rm -f bin/forte

test:
	go test -timeout 30s ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

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

example-php-deployment: build
	@echo "\033[1;33mBuilding image...\033[0m"
	docker build --quiet -f examples/php/Dockerfile -t forte-example-php .
	docker run --read-only --rm -d --name forte-example-php -p 8000:8000 forte-example-php
	curl -fs --retry 5 --retry-delay 2  http://0.0.0.0:8000/ | grep -q '"status":"ok"' && \
		(echo "\033[0;32m✓ PHP deployed app health check passed\033[0m"; docker stop forte-example-php) || \
		(docker logs forte-example-php; docker stop forte-example-php; exit 1)
