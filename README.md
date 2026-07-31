# go-bootstrap

> TODO: Replace this text with a short description of the problem this project solves and who it
> serves.

## Status

> TODO: State whether the project is experimental, under active development, production-ready, or
> deprecated. Link to a roadmap or issue tracker when one exists.

## Features

- HTTP server with graceful shutdown and configurable timeouts.
- Environment-based configuration loaded once at startup.
- Structured, level-based logging.
- Health endpoint and panic-recovery middleware.
- Optional MySQL and DynamoDB Local development dependencies.
- Unit tests, race-detector support, container builds, and clean source packaging.
- TODO: Replace or extend this list with project-specific capabilities.

## Architecture

Application code is grouped by responsibility:

- `handler` owns HTTP transport concerns and route registration.
- `service` owns business rules and use cases.
- `repository` owns persistence implementations.
- `client` owns outbound service or queue integrations.
- `internal/pkg` contains reusable infrastructure that is independent of application logic.

Keep handlers thin and move business decisions into services. Add interfaces only where they make
dependencies or tests clearer.

### Project structure

```text
cmd/server/main.go                  executable entrypoint
config/                             startup configuration
Dockerfile                           production-style image build
Makefile                             development commands
internal/
  go-bootstrap/
    handler/                        HTTP handlers and route registration
    client/                         outbound integrations
    service/                        business logic
    repository/                     persistence implementations
  pkg/
    healthcheck/                     process health handler
    httpserver/                     graceful HTTP server wrapper
    logger/                         process-wide structured logger
    middleware/                     reusable HTTP middleware
infra/
  local/
    compose.dependencies.yaml        optional local dependencies
    .env.example                     local dependency configuration template
scripts/
  package.sh                         clean source archive builder
```

<!-- init-project:start -->
## Initialize a project

Copy this repository into the project directory and run the initializer once:

```bash
./scripts/init-project.sh <project-name> [module-path]
```

For example:

```bash
./scripts/init-project.sh payments-service github.com/acme/payments-service
```

The module path defaults to the project name. The initializer updates the Go module, imports,
application directory, binary and archive names, local Compose project, and documentation. It then
formats the renamed Go files and removes itself. It does not rename the enclosing directory.
<!-- init-project:end -->

## Getting started

### Prerequisites

- Go 1.26 or newer.
- GNU Make.
- Docker with Compose support when local dependencies or a container image are needed.

### Install dependencies

```bash
go mod download
```

### Run locally

```bash
make check
make run
```

The server listens on <http://localhost:8083> by default.

## Configuration

Configuration is read from environment variables once at startup. Restart the process after
changing a value.

| Variable | Default | Purpose |
| --- | --- | --- |
| `HTTP_ADDRESS` | `:8083` | HTTP listener address |
| `HTTP_READ_HEADER_TIMEOUT` | `5s` | Maximum time to read HTTP headers |
| `HTTP_WRITE_TIMEOUT` | `10s` | Maximum time to write a response |
| `HTTP_IDLE_TIMEOUT` | `60s` | Maximum idle keep-alive time |
| `HTTP_SHUTDOWN_TIMEOUT` | `10s` | Graceful shutdown deadline |
| `LOG_LEVEL` | `info` | Application log level |

> TODO: Document every project-specific variable, whether it is required, its default, and how its
> secret value is supplied in each environment.

## API

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/healthz` | Process health check |
| `GET` | `/users` | Example endpoint; replace or remove it |

> TODO: Document authentication, request and response formats, error conventions, pagination, rate
> limits, and a link to the API specification when applicable.

## Local dependencies

MySQL and DynamoDB Local are available for development but are not required by the example API.
Start only the dependencies the project uses.

```bash
make deps-up
make deps-ps
```

The first command creates `infra/local/.env` from `infra/local/.env.example`. Edit the local file to
change ports or credentials.

Default connections:

```text
MySQL:    127.0.0.1:3306
          database=db user=user password=password
DynamoDB: http://127.0.0.1:8000
          region=us-east-1 credentials=local/local
```

Stop the dependencies without deleting their data:

```bash
make deps-down
```

> TODO: Remove unused dependencies and document required schemas, migrations, seed data, and
> external services.

## Development

| Command | Description |
| --- | --- |
| `make fmt` | Format Go source files |
| `make check` | Check formatting, run `go vet`, and run tests |
| `make race` | Run tests with the race detector |
| `make build` | Build `bin/go-bootstrap` |
| `make run` | Run the API locally |
| `make deps-up` | Start optional local dependencies |
| `make deps-down` | Stop optional local dependencies |
| `make docker-build` | Build the local container image |
| `make zip` | Create `dist/go-bootstrap.zip` |
| `make clean` | Remove generated binary and archive files |

### Testing

Keep focused tests beside the code they cover. Run `make check` before submitting a change and
`make race` for changes involving shared state, concurrency, server lifecycle, or middleware.

> TODO: Document integration, contract, end-to-end, load, and test-data workflows as they are added.

## Observability

- Application logs are written as structured output and controlled by `LOG_LEVEL`.
- `GET /healthz` reports process health.

> TODO: Document log fields, metrics, traces, dashboards, alerts, service-level objectives, and
> operational runbooks.

## Deployment

Build the container image with:

```bash
make docker-build
```

The image runs the server binary and expects runtime configuration through environment variables.

> TODO: Document environments, deployment commands, infrastructure ownership, database migration
> order, health gates, rollback steps, and release approval requirements.

## Packaging

Run `make zip` to create `dist/go-bootstrap.zip`. The source archive excludes Git metadata, local
editor configuration such as `.vscode`, environment files, build outputs, coverage reports, test
profiles, and operating-system metadata.

## Security

- Never commit credentials, tokens, private keys, or local environment files.
- Keep dependencies updated and review new dependencies before adding them.
- Validate untrusted input at the system boundary.

> TODO: Document authentication, authorization, data classification, secret management, dependency
> scanning, and the private channel for reporting vulnerabilities.

## Known limitations

> TODO: List important constraints, unsupported use cases, scaling limits, and accepted technical
> debt. Remove this section when there are no known limitations.

## Continuous integration

The `Check` workflow in `.github/workflows/check.yml` runs `make check` for every pull request
targeting `master` and every push to `master`. Update its branch filters if the repository uses a
different default branch.

> TODO: Add integration tests, security checks, release jobs, and deployment gates as the project
> evolves.

## Contributing

Keep changes focused, explain important assumptions, add tests for changed behavior, and run
`make check` before opening a review.

> TODO: Add the branching strategy, commit conventions, review requirements, and links to team
> engineering standards.

## Maintainers

> TODO: Add the owning team, support channel, escalation contact, and on-call rotation.

## License

> TODO: State the project license and add the corresponding license file. For private projects,
> describe the internal usage restrictions instead.
