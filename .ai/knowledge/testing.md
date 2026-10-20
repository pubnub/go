# Testing

Where tests live and how to match existing style. What to write for a given change is specified in the task; match nearby tests for the same feature.

## Layout

- **Root unit tests** (`*_test.go`, `package pubnub`): same package as production code. Dominant assert library: `github.com/stretchr/testify/assert`. Common fixtures: `NewDemoConfig()`, path helpers from `tests/helpers` (`AssertPathsEqual`, `AssertQueriesEqual`).
- **E2E** (`tests/e2e`, `package e2e`): loads `../../.env` via godotenv when present; env keys `PUBLISH_KEY` / `SUBSCRIBE_KEY` and PAM variants. Live DataSync tests use a dedicated PAM keyset via `DS_PUBLISH_KEY` / `DS_SUBSCRIBE_KEY` / `DS_SECRET_KEY` and skip when any of those is unset — remaining e2e tests keep the original keys. Stubbing via `tests/stubs` interceptor or custom `RoundTripper`s.
- **Example tests** under `examples/snippets/api/` are executed as part of the full suite.

## Commands

Confirm current targets in the `Makefile` and CI workflow. Typical goals:

| Goal | Command |
|------|---------|
| Full suite (CI shape) | `make test` → `scripts/run-tests-gotestsum.sh` |
| CI exact invocation | `bash ./scripts/run-tests-gotestsum.sh -mod=mod testname` |
| Root package only | `make test-quick` |
| Fast short run | `make test-dev` |
| E2E only | `make test-integration` |
| Utils / helpers | `make test-utils` / `make test-helpers` |

Optional local tools: `make install-tools` (gotestsum, gocovmerge).

`scripts/run-tests.sh` is legacy and unused — do not use it.

When verifying a change, prefer the narrowest relevant target first, then broaden if the task requires it. Do not invent E2E keys; use demo/stub patterns when keys are unavailable.
