# Go Application Template

Template repository for HTTP services written in Go and intended to run in Docker/Kubernetes.

Replace this paragraph with the project-specific description: what the service owns, which external systems it talks to, and which API or background jobs it exposes.

## What Is Included

- HTTP API generated from OpenAPI with `oapi-codegen`.
- PostgreSQL queries generated with `sqlc`.
- Goose migrations.
- Transaction helper for PostgreSQL-backed use cases.
- Transactional outbox adapter based on `gobox`.
- JSON logging with build metadata.
- Liveness/readiness probes.
- Docker image build with version, commit hash, and build time.
- Basic lint, format, generation, and test commands.

## Requirements

- Go 1.27
- Docker
- `sqlc`
- `oapi-codegen`
- `yq`
- `golangci-lint`
- `mockgen` is run through `go run go.uber.org/mock/mockgen@v0.6.0`

## Configuration

Configuration is read from environment variables in `internal/app/config.go`.

| Variable | Default | Description |
| --- | --- | --- |
| `IMAGE_NAME` | `docker.io/fr33dman/go-template` | Docker image name used by `make build`. |
| `LOGGER_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, or `error`. Logs are always JSON. |
| `SERVER_HOST` | `0.0.0.0` | Host for API and probes servers. |
| `SERVER_API_PORT` | `8000` | HTTP API port. |
| `SERVER_HEALTHCHECK_PORT` | `8081` | Liveness/readiness probes port. |
| `SERVER_METRICS_PORT` | `9090` | Reserved metrics port. |
| `EMBEDDINGS_CONFIG_BASE64` | required | Base64-encoded YAML embedding providers and model allowlist. |
Embedding configuration is loaded by `internal/clients/embeddings` from
`EMBEDDINGS_CONFIG_BASE64`. The decoded YAML shape is:

```yaml
providers:
  - name: openai
    type: openai_compatible
    base_url: https://api.openai.com/v1
    api_key: secret
    timeout: 30s
models:
  - name: text-embedding-3-small
    provider: openai
    upstream_model: text-embedding-3-small
    dimensions: 1536
    max_input_tokens: 8191
    multilingual: true
    enabled: true
```

Only enabled models from this allowlist are exposed by `GET /models` and accepted
by `POST /generate`.

## Commands

```sh
make init
```

Initializes local files and dependencies.

```sh
make gen
```

Runs SQL and OpenAPI code generation.

```sh
make sqlcgen
```

Runs only SQL code generation from `sql/queries/sqlc.yaml`.

```sh
make test
```

Runs Go tests.

```sh
make lint
```

Runs `golangci-lint`.

```sh
make format
```

Formats Go files.

```sh
make build
```

Builds the Docker image. The image tag and binary metadata are derived from:

- `VERSION`: `git describe --tags --always`
- `COMMIT_HASH`: last commit short hash
- `DATETIME`: UTC build time

These values are passed to the image build and written into `internal/build`.

## GitHub Actions

Pull requests run tests, linters, and a Buildah image build without publishing
the image. Each operation runs as a separate GitHub Actions job. Pushing a Git
tag builds and publishes the image to Docker Hub.

Configure the following repository settings before using the release workflow:

- repository variable `DOCKERHUB_REPOSITORY`, for example `my-user/my-service`;
- repository secret `DOCKERHUB_USERNAME`;
- repository secret `DOCKERHUB_TOKEN` containing a Docker Hub access token.

The release workflow publishes the image with the exact Git tag as its Docker
tag.

The workflows use `go-runners` for Go jobs and `buildah-runners` for image
jobs. These labels should point to separate ARC runner scale sets in
Kubernetes. The Buildah runner must support GitHub Actions container jobs and
provide the permissions required by Buildah, commonly through a privileged
runner Pod or a configured rootless Buildah setup.

## API

The API contract is defined in `api/openapi.yml`.

Generated files:

- `internal/server/gen/server.gen.go`
- `internal/server/gen/spec.gen.go`
- `internal/server/gen/types.gen.go`

After editing OpenAPI, run:

```sh
make gen
```

## Database

Migrations live in `sql/migrations`.

SQL queries live under `sql/queries/<repository>/`.

The sqlc configuration is `sql/queries/sqlc.yaml`. Generated repository code is written to `internal/repo/<repository>/sqlc`.

After editing SQL queries or sqlc config, run:

```sh
make sqlcgen
```

## Runtime

The main API server is started by `cmd/server`.

The migration binary is started by `cmd/migrate`.

The application handles `SIGINT` and `SIGTERM`, gracefully shuts down API and probes servers, and closes the PostgreSQL pool.

Logs are emitted as JSON and include the application version.

## Kubernetes

Kubernetes manifests are under `deploy/` and use the `go-template` namespace.
Before applying them, provide a `go-template-postgres` Secret in that
namespace with the keys `host`, `port`, `database`, `username`, and `password`.

The manifests use `docker.io/fr33dman/go-template:latest` as the template
image. Override it in an application overlay or update the `images` section in
`deploy/kustomization.yaml` for the concrete Docker Hub repository and tag.

Apply the base manifests with:

```sh
kubectl apply -k deploy/
```
