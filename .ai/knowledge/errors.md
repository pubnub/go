# Errors and validation

Behavior-changing rules for errors in this SDK. Inspect `pnerr/` and existing peers when implementing something new.

## Rules

- Use typed errors from `pnerr/`.
- Construct validation failures with `newValidationError(o, msg)`.
- Prefer existing `StrMissing…` / `StrInvalid…` constants in `pubnub.go` over inline strings; add a new constant next to the others when introducing a new message.
- Unit tests often assert on these exact strings — do not change message text casually.

## Construction

```go
return newValidationError(o, StrMissingSubKey)
```

Message shape (from `pnerr.NewValidationError`):

`pubnub/validation: pubnub: <OperationType.String()>: <msg>`

## Common typed errors

See `pnerr/` for the authoritative set. Typical kinds include validation, build-request, server, connection, and response-parsing errors (each with a stable `Error()` prefix such as `pubnub/validation: …`).

JSON parse failures commonly use `pnerr.NewResponseParsingError` and log via `loggerManager`.
