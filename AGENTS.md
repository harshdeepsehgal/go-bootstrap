# AI working agreement

## Implementation

- Present a plan for the implementation before making any changes for a prompt.
- Prefer the Go standard library and add dependencies only for a concrete requirement.
- Record important assumptions before adding authentication, persistence, concurrency, or queues.
- Keep changes explainable and avoid unrelated refactors.

## Package boundaries

- Follow the existing package boundaries: handlers own transport concerns, services own business
  logic, repositories own persistence, and clients own outbound integrations.
- Keep reusable infrastructure in `internal/pkg` only when it is independent of application logic.
- Keep common utility functions `internal/pkg/utils` package if they are independent of application logic.

## Unit tests

- Keep tests beside the code they cover and use Go's standard `testing` package.
- Name tests after the function or method under test, such as `TestNew` or `TestServer_Run`.
- Prefer GoLand-style table-driven tests when behavior has multiple meaningful cases. Use a
  `tests` slice with a `name` field and run each case with `t.Run`.
- Keep a single lifecycle or concurrency scenario as a direct test when a table would add noise.
- Test observable behavior rather than private implementation details. Cover success, failure, and
  boundary cases that materially affect the contract.
- Use `httptest` for HTTP handlers and middleware.

## Verification and hygiene

- Add focused tests for changed behavior and run `make check` before finishing.
- Never commit secrets, local environment files, editor state, or generated build artifacts.
