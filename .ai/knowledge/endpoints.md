# HTTP endpoints

How request endpoints are implemented in this SDK. Prefer existing peers in the same feature area over inventing structure. Confirm details against current code when they matter.

## Core pattern

All request endpoints implement the `endpoint` interface in `endpoints.go`. Feature opts **embed** `endpointOpts` and override what they need.

Required / commonly overridden methods:

- `validate() error`
- `buildPath() (string, error)`
- `buildQuery() (*url.Values, error)`
- `buildBody() ([]byte, error)` when the request has a body
- `buildHeaders() (map[string]string, error)` when non-default headers are needed
- `httpMethod() string` — defaults to `"GET"` on `endpointOpts`; override for `POST`, `POSTFORM`, `PUT`, `PATCH`, `DELETE`
- `operationType() OperationType`

Typical shape (see any existing endpoint file):

1. Unexported `*…Builder` holding `opts *…Opts`
2. `new…Builder(pn)` → `new…BuilderWithContext(pn, pn.ctx)`
3. Fluent setters returning the builder
4. `Execute() (response, StatusResponse, error)` → usually `LogUserInput` → `executeRequest(opts)` → response parser

`executeRequest` in `request.go` validates, builds the signed URL, issues the HTTP call (including `PUT` / `PATCH`), and returns raw body + `StatusResponse`.

## Entry points and shared helpers

Public entry points are either methods on `*PubNub` or on a feature service (e.g. `DataSyncService` via `pn.DataSync`). Many APIs also expose a `…WithContext(ctx Context)` twin. Follow the peers for the feature you are changing.

Shared URL helpers live in `endpoints.go`: `defaultQuery`, `SetQueryParam`, `SetArrayTypeQueryParam`, `buildURL`, `newValidationError`.

Feature-shared constants and domain types often live in `*_common.go` (e.g. `datasync_common.go`, `objects_common.go`, `files_common.go`).

## Operations and enums

- Operation enums are `PN…Operation` values of type `OperationType` in `enums.go`.
- When adding an operation, add both the iota constant **and** a matching `case` in `OperationType.String()`. Validation errors use that string via `newValidationError`.
- Comments in `enums.go` state that **enum order must be maintained** for Objects API sections — do not reorder existing Objects constants.
- Telemetry short codes for new operation families live in `telemetry_manager.go`; follow nearby DataSync / Objects cases.
