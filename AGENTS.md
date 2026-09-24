# Repository Instructions

This repository is a Go service template. Keep changes small, generated-code aware, and consistent with the existing layering.

## Structure

- `cmd/server`: API service entrypoint.
- `cmd/migrate`: migration entrypoint.
- `internal/app`: application bootstrap, config loading, dependency wiring, resource cleanup.
- `internal/build`: build-time metadata injected by Docker build flags.
- `internal/core`: domain entities, value objects, and domain constants.
- `internal/errors`: application-level sentinel errors, package name `apperrors`.
- `internal/services`: use cases and service-owned ports.
- `internal/repo/<name>`: repository adapters.
- `internal/repo/<name>/sqlc`: generated sqlc code for one repository/query group.
- `internal/outbox`: outbox adapter implementation.
- `internal/server`: HTTP server setup, middleware, OpenAPI validation, error mapping.
- `internal/server/handlers`: HTTP handlers that map generated DTOs to service calls.
- `internal/server/gen`: generated OpenAPI code.
- `internal/heartbeat`: app-specific probe wiring and probe server.
- `internal/logger`: JSON `slog` setup.
- `pkg/database`: reusable database helpers, including transaction context support.
- `pkg/probes`: reusable probe state and heartbeat loop.
- `api/openapi.yml`: OpenAPI contract.
- `api/oapi-codegen`: OpenAPI generation configs.
- `sql/migrations`: goose migrations.
- `sql/queries`: sqlc configs and raw SQL queries.

## Generated Code

Do not edit generated files by hand:

- `internal/server/gen/*.go`
- `internal/repo/*/sqlc/*.go`
- `internal/services/mocks/*.go`

Use:

```sh
make gen
```

for all generation, or:

```sh
make sqlcgen
```

for SQL generation only.

## Adding A New HTTP Endpoint

1. Add or change the operation in `api/openapi.yml`.
2. Run `make gen`.
3. Add a method to `internal/server/handler.go` if a new generated operation must be routed to a feature handler.
4. Add or update a handler in `internal/server/handlers`.
5. Map generated request/response DTOs at the HTTP boundary. Do not pass generated OpenAPI types into services.
6. Add or update a service method in `internal/services` when the operation has use-case logic.
7. Wire the handler/service in `internal/app/container.go` and `internal/server/api.go` if a new dependency is introduced.
8. Add focused tests for validation/error mapping or service behavior when there is non-trivial logic.

OpenAPI request validation is enabled in `internal/server/api.go` through `gen.GetSwagger()` and `echomiddleware.OapiRequestValidator`.

## Adding A Database Query

1. Add SQL under `sql/queries/<repository>/`.
2. If this is a new repository group, add a new entry to `sql/queries/sqlc.yaml`.
3. Run `make sqlcgen`.
4. Add a repository method in `internal/repo/<repository>`.
5. Map sqlc models to `internal/core` models inside the repository package.
6. Extend the service-owned repository interface in `internal/services` only with methods the service actually needs.
7. Wire the concrete repository in `internal/app/container.go`.

Repositories should use `database.DBTX(ctx, fallback)` so calls inside `Transactor.WithinTx` use the active transaction and calls outside a transaction use the pool.

## Transactions

Services depend on this port:

```go
type Transactor interface {
    WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

Use it when a use case must commit several operations atomically:

```go
err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
    // repository and outbox calls here use the transaction from ctx
    return nil
})
```

Do not pass `pgx.Tx`, `pgxpool.Pool`, or sqlc types into services. PostgreSQL details belong in repository/outbox/database adapters.

## Errors

Use `internal/errors` (`package apperrors`) for application-level sentinel errors such as:

- `ErrNotFound`
- `ErrInvalidArgument`

Use `internal/core` only for true domain errors and invariants. HTTP status mapping belongs in `internal/server/api.go`.

## Configuration

Add environment-backed settings in `internal/app/config.go`.

Keep config grouped by responsibility:

- `PostgresConfig`
- `ServerConfig`
- `LoggerConfig`
- new groups for new subsystems

Update `.env.template` and `README.md` when adding a new required or user-facing variable.

## Logging

Use `internal/logger`.

Logs must stay JSON. Do not add plain `log.Printf`/`fmt.Println` runtime logging. Request logging is configured in `internal/server/api.go`.

Build metadata lives in `internal/build` and is injected by Docker `-ldflags`.

## CI/CD

GitHub Actions are defined in `.github/workflows`:

- `pull_request.yml` runs tests, linters, and builds the image with Buildah in separate jobs without pushing it;
- `release.yml` builds and pushes an image on every Git tag in separate jobs.

The release workflow expects the `DOCKERHUB_REPOSITORY` repository variable and
the `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` repository secrets.

The workflows use `go-runners` and `buildah-runners` as runner labels. These
labels are expected to be provided by separate ARC runner scale sets.

## Kubernetes Manifests

The base manifests in `deploy/` use the `go-template` namespace and expect a
`go-template-postgres` Secret with `host`, `port`, `database`, `username`, and
`password` keys. Update the image through a Kustomize overlay or the
`images` section in `deploy/kustomization.yaml` for a concrete application.

## Probes

Probe state is created in `internal/app/container.go` and passed explicitly to the probes server and heartbeat loop.

Do not reintroduce global probe state.

Readiness should represent external dependencies required to serve traffic. Liveness should represent the process being alive; fatal server errors should return from `RunServer` and terminate the process.

## Outbox

Services know only the `EventPublisher` port and `services.Event`.

The `gobox` implementation belongs in `internal/outbox`. Outbox publication should happen inside `Transactor.WithinTx` when it must be atomic with database writes.

## Tests

Keep tests focused on meaningful behavior:

- service transaction/event behavior;
- HTTP validation and error mapping;
- small utility behavior with actual logic.

Use `mockgen` for service port mocks. Regenerate mocks with `make gen`.

## Before Finishing

Run:

```sh
make format
make test
make lint
```

If a generated file changes, make sure the source contract/config that generated it is also updated.
