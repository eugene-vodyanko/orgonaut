# Development Guide

Contains useful information for developers.

## Project versioning
- Semantic versioning is used, examples: `0.0.1-alpha`, `0.2.1-beta`, `1.0.0-rc`.
- Instead of a separate version file, the use of `git tags`.

## TODO List
- Metrics in Prometheus.
- Integration tests (mock code generation).
- Benchmarks.
- CI.
- Any primary key support.
- Fine grained error handling.

## FAQ

### Launching the local dev. environment in docker-compose:
Deploy and launch a container with `Kafka`, `Kafka-UI`:
```shell
docker-compose up -d
```

Stop and clean:
```shell
docker-compose down -v
```

Create Kafka topic with two partitions with cli:
```sh
docker exec -it kafka kafka-topics --create --topic messages --bootstrap-server localhost:29092 --partitions 2
```

Local Kafka UI:
http://localhost:8082/

### Build and launch the service
Directly:
```shell
go mod tidy                      
go run cmd/app/main.go
```
Or the same thing with `make`:
```shell
make run
```

### Linter
Run a static analysis using `golangci-lint`:
```shell
make lint
```