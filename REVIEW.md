# REVIEW.md

Review instructions for the PubNub Go SDK. The review action provides the general procedure, the severity definitions, the finding format with examples, and the handling of earlier findings. This file only holds what is specific to this repository.

## Scope

- Apply `AGENTS.md`. When the diff touches an area below, read the matching note and apply it. Name the rule in each finding.
  - HTTP endpoints — `.ai/knowledge/endpoints.md`
  - Errors and validation — `.ai/knowledge/errors.md`
  - Tests — `.ai/knowledge/testing.md`
  - DataSync — `.ai/knowledge/datasync.md`
- Paths in these instructions are relative to the repository root, which is your working directory.
- Code wins when a doc and the implementation disagree. Cite the code.
- Trace production changes from the public entry point through validation, request construction, decoding, and tests.
- Search for related declarations, helpers, fixtures, snippets, and release metadata outside the diff.

## Writing findings

Write the way a helpful teammate would in a code review: plain words, say what you noticed and why it matters, then suggest a fix. Explain, don't just assert.

Flag an issue only when this diff introduced it, it is concrete, and the author would fix it. Skip speculative concerns and pre-existing problems.

## Review checklist

### Public API

The module is `github.com/pubnub/go/v10`. The public surface is the flat `pubnub` package, `crypto/`, and `pnerr/`. Callers use methods on `*PubNub`, or `pn.DataSync` (`DataSyncService`) for DataSync. Builders and opts are unexported. Fluent setters and `Execute` are the caller-facing surface.

- Trace every parameter of a new or changed method to the correct request field. None may be dropped or swapped.
- Keep the `*PubNub` method thin. Validation, path, query, body, and headers belong on the endpoint opts, as the nearest peer does.
- Exported identifiers use Go-idiomatic acronyms: `UUID`, `HTTP`, `TLS`, `TTL`. Public response and domain types are often prefixed `PN`.
- New example or production code uses `NewConfigWithUserId`. A new call to deprecated `NewConfig` is `minor`.
- Keep declarations unexported unless consumers call or implement them. Widening access needs a concrete use case in the PR description and matching tests.
- A change to an exported signature, type, or observable behavior on `*PubNub`, a `PN…` type, or `crypto` that makes callers fail to compile, or that changes the meaning, is `major`.

### Breaking changes

Breaking changes may be intentional, but must be declared.

- API break: removal, renaming, or retyping of an exported symbol; a changed result shape; a parameter added, removed, or reordered; a changed method set on an exported interface.
- Behavior break: a changed default, limit, ordering guarantee, error text, or callback frequency.
- The release bot owns `Version` in `pubnub.go` and `version` / `changelog` in `.pubnub.yml`. Do not report a missing version bump or changelog entry.
- Report one `high` finding for all undeclared source or behavior breaks. Anchor it to the most significant, and name the others.

### Concurrency

- A new goroutine in subscribe, listeners, or request handling must stop on context cancel or client close.
- A data race on shared maps, listener state, or client teardown is `major`.
- Closing a channel, or signaling teardown, must be safe to call more than once. Peers use `sync.Once` and a single close.
- Report a user callback invoked while a lock that the callback needs is held.

### Security

- `Config.GetLogString` masks secrets: `SecretKey` through `redactPrefix` (first 8 characters and the length; shorter values become `***`), and `AuthKey` and `CipherKey` fully masked. New logs of keys, tokens, cipher keys, authorization headers, or signed URLs must follow that.
- No real keys or tokens in the diff, including tests, snippets, and examples. Tests use `NewDemoConfig`, `tests/stubs`, or environment variables, and skip when the keys are unset.
- Crypto changes in `crypto/` preserve nonce uniqueness, authentication checks, ciphertext layout, and compatibility with existing payloads. Decryption failures stay generic. Do not recommend replacing a compatibility cryptor unless the change makes it less safe.
- A new direct dependency in `go.mod` needs a reason in the PR description.

### Build and release integrity

CI does not catch these.

- Minimum Go is 1.25 (`go 1.25.0` in `go.mod`, Go 1.25 in `.github/workflows/`). Report a toolchain or language feature newer than that.
- A hand edit of release metadata is `minor`: `Version` in `pubnub.go`, or `version` / `changelog` in `.pubnub.yml`.
- Module path stays `github.com/pubnub/go/v10`.

### Networking

- Changed endpoints keep the `endpoint` interface and the nearest builder: `validate`, `buildPath`, `buildQuery`, `buildBody`, `buildHeaders`, `httpMethod`, `operationType`, then `Execute`.
- HTTP method, path, percent-encoding, query names and omission rules, headers, body encoding, and signing inputs match the task and in-repo peers. A divergence is `major`.
- A new `PN…Operation` updates both the iota constant and `OperationType.String()`. Reordering existing Objects constants in `enums.go` is `major` (`ENUM ORDER needs to be maintained`). A new operation family also needs the telemetry case peers have in `telemetry_manager.go`.
- Decoding handles documented success shapes and maps status codes, service errors, and malformed responses without dropping useful error information. JSON parse failures use `pnerr.NewResponseParsingError`.
- Pagination cursors, limits, and offsets travel in the correct direction.

### DataSync

Confirm against `.ai/knowledge/datasync.md` and `datasync_*.go`.

- Public methods are `Remove*`. Files and unexported builders stay `delete*`.
- Entity and Relationship operation constants have no `DataSync` infix. User, Channel, and Membership constants do. Renaming or reordering either group is `major`.
- `UUIDs()` / `UUIDsPattern()` do not authorize DataSync User. User CRUD is `Users()` / `UsersPattern()`.
- Set is a full replacement. Update is RFC 6902 JSON Patch. `Filter` and `FilterFast` are mutually exclusive. An unset class version must not be sent as `0`.
- Subscribe message type 5 with `source: "data-sync"` and `version: "1.0"` is a DataSync event. Other versions are ignored.

### Errors

- Errors that escape the package use `pnerr/`.
- Validation failures use `newValidationError` and the `StrMissing…` / `StrInvalid…` constants in `pubnub.go`. The resulting text is `pubnub/validation: pubnub: <OperationType.String()>: <msg>`. Changing that text changes unit asserts. An unexplained change is `major`.

### Tests

Test layout is in `.ai/knowledge/testing.md`. Enforce only the rules below.

- A bug fix or behavior change comes with a test that would fail without it. A miss is `major`. Comments, formatting, renames, and generated files do not need a new test.
- Assert the behavior that changed: path, query, body, validation error text, response parse, crypto, or DataSync event handling.
- Root unit tests are `*_test.go` in package `pubnub`, usually `github.com/stretchr/testify/assert`, with `tests/helpers` for paths and queries and `tests/stubs` for HTTP.
- A stub unit test is enough when it covers the behavior. Do not require a live end-to-end test on top of it. `tests/e2e` skips when `PUBLISH_KEY` / `SUBSCRIBE_KEY` (and the PAM or `DS_*` variants) are unset.
- Do not report that tests were not executed.
- `scripts/run-tests.sh` is legacy. CI uses `scripts/run-tests-gotestsum.sh`.

### Snippets

`examples/snippets/api/` example tests compile and run in CI. Check that they teach the public API correctly.

- Use the public `pubnub` API and `NewConfigWithUserId`. No real secret keys in client code. PAM examples keep the secret on the backend or behind the existing env helper.
- Preserve `// snippet.` and `// snippet.end` markers and IDs. `// snippet.hide` / `// snippet.show` wrap compilation scaffolding only.
- Add a snippet for new public API in the matching `examples/snippets/api/` file, modeled on the nearest snippet, when that file is the consumer-facing sample for the area.

### Documentation

- `README.md` is the consumer quickstart. An update is required when the diff changes a fact it states or a public step it shows: Go version, module path, or the quickstart calls. A new endpoint that follows existing peers does not need a new paragraph.
- `AGENTS.md` and `.ai/knowledge/` need an update when the diff changes a rule one of those files states. A new endpoint that follows existing peers does not need a new paragraph.
- Godoc claims match the code: parameter names, defaults, `nil` and empty-value behavior. Ordering, retries, persistence, or delivery are promised only when source and tests establish them.
- Wording and punctuation problems in shipped docs are `nitpick`.

## Do not report

- Formatting that `gofmt` covers, import grouping, and line length.
- Taste and vague suggestions: naming preferences, extra abstractions, and performance concerns without a concrete hot path.
- Missing version bumps and changelog edits. The release pipeline owns them.
- Build and test failures already reported by required CI.
- Pre-existing issues the diff did not introduce.
