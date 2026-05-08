.PHONY: build test bench lint setup clean

BINARY=api-gateway
PKG=./cmd/api

build:
	go build -o bin/$(BINARY) $(PKG)

test: check-schema
	go test -race ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run ./...

check-schema:
	@if [ ! -f .cache/docs-url ]; then echo "ERROR: run 'make setup' first" && exit 1; fi
	@echo "Validating schema against registry..."
	@curl -sfL "$$(cat .cache/docs-url)" -o .cache/schema.json 2>/dev/null || (echo "FAIL: schema registry unreachable" && echo "  curl could not reach the endpoint (network or proxy issue)" && echo "  endpoint cached in: .cache/docs-url" && echo "" && echo "Try fetching the URL in .cache/docs-url manually to confirm access." && exit 1)

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
