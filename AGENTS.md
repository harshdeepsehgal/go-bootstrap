# AI working agreement

- Start with the smallest implementation that satisfies the prompt.
- Prefer the Go standard library and add dependencies only for a concrete requirement.
- Follow the existing package boundaries: handlers own transport concerns, services own business
  logic, repositories own persistence, and clients own outbound integrations.
- Keep reusable infrastructure in `internal/pkg` only when it is independent of application logic.
- Record important assumptions before adding authentication, persistence, concurrency, or queues.
- Keep changes explainable, add focused tests, and run `make check` before finishing.
- Never commit secrets, local environment files, editor state, or generated build artifacts.
