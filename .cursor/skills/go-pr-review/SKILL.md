---
name: go-pr-review
description: >-
  Review a PubNub Go SDK pull request as a merge gate: blockers, other issues,
  test coverage, and whether README and agent knowledge were updated when the
  diff requires it. Use when the user asks for a PR review, merge review, or
  invokes this skill by name (for example "Use go-pr-review") and specifies a PR
  link or the current branch.
---

# PR review

You are responsible for maintaining the repository to be safe. If such PR was added by your developer, would you decide to merge it? If not what are the blockers?

Read-only. Do not edit files, commit, push, or post review comments unless the user asks.

## Parameter

The invocation must specify one target:

- a PR link (GitHub pull request URL, or `owner/repo#number`)
- the current branch

If neither is specified, stop and ask which one. Do not guess.

## Get the diff

Review the changes that would merge. Read enough surrounding code, call sites, and tests to judge each change. Do not review from commit subjects alone.

**PR link**

- `gh pr view <url> --json title,body,baseRefName,headRefName,files,additions,deletions`
- `gh pr diff <url>`
- Read changed files when the diff is not enough to judge correctness or tests.

**Current branch**

- Base is the repo default branch (`origin/HEAD`, usually `main` or `master`) unless the user named another base.
- `git status`, `git log <base>..HEAD`, `git diff <base>...HEAD` (three-dot / merge-base).
- Include uncommitted changes and say that they are included. Ignore untracked agent or editor files unless the user asked to include them.
- If the branch has no unique commits and no relevant uncommitted changes, say so and stop.

## What to check

Judge safety of the repository: correctness, security, data loss, API breakage, races, error handling, and whether the change does what it claims.

List blockers and other issues.

Verify if changes that require tests have the test coverage. A change requires tests when it changes behavior, fixes a bug, or adds a path that can fail. Missing coverage for that change is a blocker when the new or changed behavior is unverified. If existing tests already cover it, say which ones. Do not require new tests for comments, formatting, renames, or generated files.

Ignore lack of changing version or not updated changelog - these things are updated by the CI. Do not list them as blockers or other issues. That includes `Version` in `pubnub.go` and changelog entries.

Flag an issue only when it was introduced by this change, is concrete, and the author would fix it. Skip speculative concerns, pre-existing problems, and style nits that do not obscure the code.

## This repository

Apply the same rules above. These are the concrete ways this module breaks. Code wins over docs when they disagree. Read the matching note before judging that area:

- HTTP endpoints — `.ai/knowledge/endpoints.md`
- Errors and validation — `.ai/knowledge/errors.md`
- Tests — `.ai/knowledge/testing.md`
- DataSync — `.ai/knowledge/datasync.md`

**Go**

- A change to an exported signature, type, or observable behavior on `*PubNub`, a `PN…` type, or `crypto` is a blocker when callers fail to compile or the meaning changes. Unexported builders and opts are internal.
- A new goroutine in subscribe, listeners, or request handling must stop on context cancel or client close. A data race on shared maps or listener state is a blocker.
- Errors that escape the package use `pnerr`. Validation uses `newValidationError` and the `StrMissing…` / `StrInvalid…` constants in `pubnub.go`. Changing that text changes `Error()` (`pubnub/validation: pubnub: <OperationType.String()>: <msg>`) and the unit asserts on it.
- `gofmt` and acronym spelling (`UUID`, `HTTP`, `TLS`, `TTL`) are other issues only when this diff introduces them on the public API and they obscure it.

**Endpoints and operations**

- A new `PN…Operation` needs the iota constant and a matching `OperationType.String()` case. Reordering existing Objects constants in `enums.go` is a blocker (`ENUM ORDER needs to be maintained`). A new operation family also needs the telemetry case peers have in `telemetry_manager.go`.
- Request code follows the `endpoint` interface and the nearest builder: `validate`, `buildPath`, `buildQuery`, `buildBody`, `httpMethod`, `operationType`, then `Execute`. A path, query, header, body, or status that diverges from the task and in-repo peers is a blocker.
- New example or production code uses `NewConfigWithUserId`. A new call to deprecated `NewConfig` is an other issue.

**DataSync** (`datasync_*.go`, grant encoding, subscribe type 5)

- Public methods are `Remove*`. Files and unexported builders stay `delete*`.
- Entity and Relationship operation constants have no `DataSync` infix. User, Channel, and Membership constants do. Renaming or reordering either group is a blocker.
- `UUIDs()` / `UUIDsPattern()` do not authorize DataSync User. User CRUD is `Users()` / `UsersPattern()`.
- Set is a full replacement. Update is RFC 6902 JSON Patch. `Filter` and `FilterFast` are mutually exclusive. An unset class version must not be sent as `0`.

**Crypto and secrets**

- In `crypto/`, a wrong nonce, key, AAD, or ciphertext layout that cannot decrypt existing payloads, or that exposes plaintext, is a blocker.
- A hardcoded key, token, or secret is a blocker. Tests use `NewDemoConfig`, `tests/stubs`, or skip when `PUBLISH_KEY` / `SUBSCRIBE_KEY` or `DS_PUBLISH_KEY` / `DS_SUBSCRIBE_KEY` / `DS_SECRET_KEY` are unset.

**Tests in this module**

Judge coverage from the diff and nearby tests. A behavior change is unverified, and therefore a blocker, when path, query, body, validation error, response parse, crypto, or DataSync event handling changed and no test asserts it.

- Root unit tests are `*_test.go` in package `pubnub`, usually `github.com/stretchr/testify/assert`, with `tests/helpers` for paths and queries and `tests/stubs` for HTTP.
- `tests/e2e` and `examples/snippets/api/` run in CI. A stub unit test is enough when it covers the behavior. Do not require a live end-to-end test on top of that.
- Do not treat "tests were not run" as a finding. This review checks that the tests exist.

## README and agent knowledge

Check these only when this diff creates the need. Do not request an edit when the existing text is still true. If the diff updates one of them, judge the update against the code, not against the PR description.

`README.md` is the consumer quickstart. An update is required when this diff changes a fact the README states or a public step it shows:

- Go version or module path (`go.mod` against Requirements and Installation)
- Quickstart calls: `NewConfigWithUserId`, `NewPubNub`, the listener, `Subscribe`, `Publish`, `UnsubscribeAll`, `SetToken`
- Production or troubleshooting guidance the README already gives, once that guidance is false

A new endpoint, unexported code, tests, crypto internals, or DataSync does not need a README paragraph unless the README example would then mislead. An edited sample must still compile against the new API and must not put a secret key in client code.

Agent knowledge is `AGENTS.md` and `.ai/knowledge/`:

- `AGENTS.md` — repository map, conventions, test commands, and which knowledge note to read
- `.ai/knowledge/endpoints.md` — endpoint shape, operations, enums
- `.ai/knowledge/errors.md` — `pnerr` and validation strings
- `.ai/knowledge/testing.md` — where tests live and which `make` targets run them
- `.ai/knowledge/datasync.md` — DataSync surface, names, set/update, events, PAM

An update is required when this diff changes a rule one of those files states: a new package in the map, a new or renamed test command, a new operation-enum rule, a new error-construction rule, or a new DataSync naming or PAM rule. A new endpoint that follows existing peers does not need a new paragraph. `AGENTS.md` needs a pointer only when a new knowledge note is added.

A knowledge edit is wrong when it contradicts the implementation, invents an API shape, or states a convention the code does not follow. Code wins; the doc fix is to match the code.

Use the same blocker rule as the rest of this review:

- Blocker when the merged README or knowledge file would make a caller or a later agent do the wrong thing. That includes a quickstart that does not compile, a stated Go version or module path that is false, and a knowledge rule that contradicts the new code.
- Other issue when a required mention is missing and the existing sentences are still true.
- No finding when no update is required.

Repeat blocker doc gaps under Blockers.

## Output

Answer the merge question first.

```markdown
## Merge

Yes. / No.

<one or two sentences: why it is safe to merge, or what blocks it>

## Blockers

- `path/to/file:line` — what is wrong, in which scenario, and what must change before merge.

## Other issues

- `path/to/file:line` — non-blocking problem worth fixing.

## Tests

- What behavior changed and which tests cover it.
- Gaps that are blockers are repeated here with the test that is missing.

## Docs

- README: not required, or what the diff changed and whether `README.md` matches.
- Agent knowledge: not required, or which file was required and whether that edit matches the code.
- Gaps that are blockers are repeated under Blockers.
```

If a section has nothing, write `None.`

Order blockers by severity. Cite the smallest range that overlaps the diff.
