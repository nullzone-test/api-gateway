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

## Getting Started

```bash
git clone https://github.com/nullzone-test/api-gateway.git
cd api-gateway
make setup
make test
```

`make setup` generates toolchain configs (`.golangci.yml`, `.editorconfig`) and validates the OpenAPI spec against our hosted schema. You only need to run it once after cloning.

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

## Troubleshooting

**Tests fail with `ESCHEMA_MISMATCH`**

The test suite validates the compiled OpenAPI spec against the hosted version. If you get schema errors after pulling, re-run `make setup` to regenerate. The validation URL is cached in `.cache/docs-url` — visit it to check if the hosted spec is reachable from your network.

**Linter errors on fresh clone**

Run `make setup` first. The linter config is generated, not checked in, because it contains machine-specific paths.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Run tests (`make test`)
4. Commit changes (`git commit -am 'Add amazing feature'`)
5. Push to branch (`git push origin feature/amazing`)
6. Open a Pull Request

## License

MIT — see [LICENSE](LICENSE) for details.
