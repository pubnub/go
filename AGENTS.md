# AGENTS.md — PubNub Go SDK

This file summarizes repository conventions for agents. It is **not** a substitute for checking the current code, `Makefile`, `go.mod`, and CI configuration. When this file conflicts with executable configuration or current code, treat the repository as authoritative and update this documentation if appropriate.

Feature-specific choices (file naming, public entry points, naming of delete/remove APIs, test scope) are defined in the task request — follow that, not this file.

## How to work in this repository

- Prefer existing peers in the same feature area over inventing structure.
- Follow task/API specifications over assumptions.
- Keep changes scoped; avoid drive-by refactors of unrelated endpoints.
- Do not invent API shapes, paths, or status semantics — follow the task, specs, and in-repo peers.
- Do not bump `Version`, rewrite changelog, or drive release metadata unless explicitly asked.
- Never invent credentials or E2E keys; use demo/stub patterns already in the tree when keys are unavailable.

## Repository map

Confirm module path / Go version from `go.mod` and CI from `.github/workflows/`. Rough layout:

- `.` — primary `pubnub` public API (flat package)
- `crypto/` — crypto module / cryptors
- `pnerr/` — typed errors
- `utils/` — URL encode, PAM helpers, UUID checks
- `tests/helpers/`, `tests/stubs/` — path/query asserts and HTTP interceptors
- `tests/e2e/` — integration tests
- `examples/snippets/api/` — example tests run in CI

SDK version constant: `Version` in `pubnub.go` (release tooling owns bumps).

## Important conventions

- Prefer `NewConfigWithUserId` over deprecated `NewConfig`.
- Exported identifiers use Go-idiomatic acronyms: `UUID`, `HTTP`, `TLS`, `TTL`.
- Public response/domain types are often prefixed `PN…`.
- Builders and opts are normally unexported; fluent setters and `Execute` are the caller-facing surface.
- Follow nearby endpoints for API shape (paths, bodies, headers, validation).
- Adding a `PN…Operation` requires updating both the iota constant **and** `OperationType.String()`.
- Do not reorder existing Objects operation constants in `enums.go`.

## Validation

Start with the narrowest relevant test, then broaden when the scope warrants it.

| Goal | Command |
|------|---------|
| Root package | `make test-quick` |
| Fast short run | `make test-dev` |
| E2E only | `make test-integration` |
| Full suite (CI shape) | `make test` |
| Utils / helpers | `make test-utils` / `make test-helpers` |

Do not use `scripts/run-tests.sh`; it is legacy. Prefer `scripts/run-tests-gotestsum.sh` (see `Makefile` / CI workflow).

## Before completing a coding task

1. Inspect the diff for unrelated changes.
2. Run the narrowest relevant tests.
3. Run broader validation when practical for the scope.
4. Report which validation commands were run and whether they passed.
5. Do not claim tests passed if they were not executed.

## Detailed repository knowledge

When modifying HTTP endpoints, read [`.ai/knowledge/endpoints.md`](.ai/knowledge/endpoints.md).

When changing errors or validation, read [`.ai/knowledge/errors.md`](.ai/knowledge/errors.md).

When adding or modifying tests, read [`.ai/knowledge/testing.md`](.ai/knowledge/testing.md).

When changing DataSync (entities, relationships, users, channels, memberships, events, or DataSync PAM), read [`.ai/knowledge/datasync.md`](.ai/knowledge/datasync.md).

Other in-repo docs: `README.md` (consumers), `DEVELOPER.md` (defaults / renames), `.pubnub.yml` (product metadata; release pipeline owns changelog/version alignment).
