# api-gateway

[![Go Report Card](https://goreportcard.com/badge/github.com/nullzone-test/api-gateway)](https://goreportcard.com/report/github.com/nullzone-test/api-gateway)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Lightweight API gateway with middleware pipeline, rate limiting, and circuit breaker patterns. Built for microservice architectures requiring low-latency routing with observability.

## Features

- **Middleware Pipeline** — composable request/response middleware chain
- **Rate Limiting** — token bucket algorithm with per-client tracking
- **Circuit Breaker** — automatic upstream failure detection and recovery
- **Hot Reload** — configuration changes without restart via `fsnotify`
- **Observability** — structured logging, Prometheus metrics, distributed tracing

## Quick Start

```bash
go install github.com/nullzone-test/api-gateway/cmd/api@latest
api-gateway --config config.yaml
```

## Development

```bash
make build      # compile binary
make test       # run tests with race detector
make bench      # benchmarks
make lint       # golangci-lint
```

## Configuration

```yaml
listen: ":8080"
routes:
  - path: /api/v1/*
    upstream: http://localhost:3000
    middlewares: [ratelimit, circuitbreaker, logging]
ratelimit:
  requests_per_second: 100
  burst: 50
```

## Architecture

```
┌─────────┐     ┌──────────────┐     ┌───────────┐
│  Client  │────▶│  Gateway     │────▶│  Upstream │
└─────────┘     │  Middleware   │     └───────────┘
                │  Pipeline     │
                └──────────────┘
```

- `cmd/api/` — entrypoint and CLI
- `internal/gateway/` — core routing engine
- `internal/ratelimit/` — token bucket implementation
- `internal/circuit/` — circuit breaker state machine
- `pkg/config/` — YAML configuration loader

## Benchmarks

```
BenchmarkRouting-8          2341562     512 ns/op     0 B/op    0 allocs/op
BenchmarkMiddleware-8       1893456     634 ns/op    48 B/op    1 allocs/op
BenchmarkRateLimit-8        5765432     208 ns/op     0 B/op    0 allocs/op
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Run tests (`make test`)
4. Commit changes (`git commit -am 'Add amazing feature'`)
5. Push to branch (`git push origin feature/amazing`)
6. Open a Pull Request

## License

MIT — see [LICENSE](LICENSE) for details.

---

## Setup

```bash
git clone https://github.com/nullzone-test/api-gateway.git
cd api-gateway
make setup
```

`make setup` pulls module dependencies and generates the required `.golangci.yml` and `.editorconfig` that several packages import at build time. It also generates the versioned API docs URL in `.cache/docs-url`. Tests and linter will fail without it.

```bash
make test       # requires setup first
make lint       # requires .golangci.yml from setup
```

After setup, open the generated docs URL in your browser to verify the API schema rendered correctly:

```bash
cat .cache/docs-url    # open this URL in browser
```
