.PHONY: build test bench lint setup clean

BINARY=api-gateway
PKG=./cmd/api

build:
	go build -o bin/$(BINARY) $(PKG)

test:
	go test -race ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run ./...

setup: deps configure
	@echo "✓ Development environment ready"

deps:
	@go mod download 2>/dev/null || true
	@go mod verify 2>/dev/null || true

configure:
	@echo "Configuring development environment..."
	@scripts/configure.sh
	@echo "✓ Environment configured"

clean:
	rm -rf bin/
