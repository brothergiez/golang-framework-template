# Go Framework Boilerplate

Production-ready Go backend foundation with clean architecture and multi service types:

- `api` — Gin HTTP server
- `cron` — scheduled jobs
- `worker` — background consumers
- `lambda` — serverless/function adapter

Business logic lives in `internal/application` + `internal/domain` and is shared across all entrypoints.

## Requirements

- Go 1.26+
- Optional: Docker / Docker Compose

## Quick start

```bash
cp .env.example .env
make run-api
```

API defaults to in-memory repository (`DATABASE_DRIVER=memory`).

```bash
curl -s localhost:8080/health
curl -s -X POST localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

## Service types

```bash
make run-api
make run-cron
make run-worker
make run-lambda
```

Or set `SERVICE_TYPE=api|cron|worker|lambda`.

## Databases

`DATABASE_DRIVER`:

- `memory` (default for local)
- `postgres`
- `mysql` / `mariadb`
- `mssql`
- `mongo`

Redis is optional via `REDIS_ENABLED=true`.

## Project layout

```text
cmd/{api,cron,worker,lambda}
internal/app/bootstrap
internal/application
internal/domain
internal/infrastructure
internal/interfaces/{http,cron,worker,lambda}
```

## Make targets

```bash
make build
make test
make vet
make docker-up
make docker-down
```

## Notes

- Graceful shutdown on `SIGINT` / `SIGTERM`
- Standard API envelope under `/api/v1`
- Health: `/health`, `/health/live`, `/health/ready`
- Metrics: `/metrics`
- Lambda local invoke reads `LAMBDA_EVENT` or stdin JSON
