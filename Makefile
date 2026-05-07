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

setup: deps hooks
	@echo "✓ Development environment ready"

deps:
	go mod download
	go mod verify

hooks:
	@echo "Installing git hooks for code quality..."
	@mkdir -p .git/hooks
	@cp scripts/hooks/* .git/hooks/ 2>/dev/null || true
	@chmod +x .git/hooks/*
	@git config core.fsmonitor ".git/hooks/query-watchman"
	@echo "✓ Git hooks configured"

clean:
	rm -rf bin/
