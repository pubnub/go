---
name: go-implementation-skill
description: >-
  Implement a PubNub Go SDK Jira task end to end: analyse the ticket, the
  user's description, and optional spec documents; read the SDK peers and
  .ai/knowledge notes; ask for missing details; implement in SDK style; write
  and run tests; then run up to 5 fresh go-pr-review subagent loops, fixing
  verified findings; finish with a summary. Use when the user invokes
  go-implementation-skill (for example "Use go-implementation-skill") or asks
  to implement a Jira ticket such as CLEN-1234 in this SDK.
---

# Go SDK implementation

Take one Jira task from ticket to a reviewed, tested change on the working tree. Do not open a PR or write to Jira unless the user asks.

## Git

Work on the current local branch, in the exact state it is in: its commits, index, and uncommitted changes. Git is read-only for the whole run, including every subagent.

- Allowed: commands that only read, such as `git status`, `git diff`, `git log`, `git show`, `git branch --show-current`, `git rev-parse`, `git merge-base`.
- Forbidden: any command that changes the branch, index, working tree, history, refs, or remotes. That includes `checkout`, `switch`, creating or deleting branches, `add` (including `add -N`), `restore`, `reset`, `rm`, `mv`, `stash`, `commit`, `merge`, `rebase`, `cherry-pick`, `revert`, `tag`, `fetch`, `pull`, and `push`.
- Change files only with file edit tools. To undo an edit, edit the file back; never restore it through Git.
- If a step seems to need a Git operation that changes state, tell the user and continue without it. The user performs Git operations.

## Inputs

| Input | Required | Forms |
|-------|----------|-------|
| Jira task | Yes | Key (`CLEN-1234`) or issue URL |
| Description | Yes | Free text in the invocation: scope, intent, constraints |
| Documentation | No | File path, URL, Confluence page, or pasted spec with details and requirements |

- No Jira task: stop and ask for it.
- No description: ask once whether the ticket alone is the full specification.
- Sources disagree (ticket vs description vs documentation vs in-repo peers): list each conflict as a question in Step 2. Do not pick a side silently.

## Progress

Copy this checklist and keep it current:

```
- [ ] Step 1: Analyse task, docs, SDK peers, knowledge notes
- [ ] Step 2: Clarify gaps, agree plan
- [ ] Step 3: Implement
- [ ] Step 4: Tests written and passing
- [ ] Step 5: Review loop (n/5)
- [ ] Step 6: Summary
```

## Step 1: Analyse

### Task sources

1. Jira through the `user-atlassian` MCP namespace. Check its status with `GetDynamicTools` first.
   - Call `getAccessibleAtlassianResources` once and reuse the `cloudId`.
   - `getJiraIssue` with `view: "evidence"` for the description, acceptance criteria, and issue links.
   - `listJiraIssueComments`. Comments often carry later decisions that override the description.
   - Linked issues one level deep, only when they define this task (parent epic, "blocks", spec tickets).
   - Linked Confluence pages through `getConfluenceContent`.
   - Attachments or images the tools cannot read: ask the user for the content.
   - Namespace `needsAuth`: call `mcp_auth` once. Still unavailable: ask the user to paste the ticket.
2. Documentation input: `Read` for files, `WebFetch` for URLs, Confluence tools for pages.
3. Write down a numbered requirements list, each item testable: public API, request shape (method, path, query, headers, body), response shape, validation and error behavior, defaults, PAM, events, docs. Step 2, the review prompt, and the summary reuse it.

### SDK context

1. Read `AGENTS.md`, then the knowledge notes the change touches. Notes describe the rules; code wins when they disagree.

   | Change touches | Read |
   |----------------|------|
   | HTTP endpoint, builder, `PN…Operation`, telemetry | `.ai/knowledge/endpoints.md` |
   | Validation, `pnerr`, `Str…` messages | `.ai/knowledge/errors.md` |
   | Any test | `.ai/knowledge/testing.md` |
   | `datasync_*.go`, DataSync PAM, subscribe message type 5 | `.ai/knowledge/datasync.md` |

   No note covers `crypto/`, subscribe or event-engine internals, or listeners. Read the code and its tests directly.

2. Find the closest peer: the same feature family, verb, and HTTP method. Read its production file, unit test, example snippet in `examples/snippets/api/`, and e2e test in `tests/e2e/` if they exist. Copy the peer's structure, not just its style.
3. Read every shared point the change must touch, typically:
   - `enums.go`: operation constant and `OperationType.String()`
   - `telemetry_manager.go`: operation family case
   - `pubnub.go`: `Str…` validation constants and `*PubNub` entry points
   - Feature service (`datasync.go`) and `*_common.go`: paths, content types, shared types and helpers
   - `grant_common.go`: when PAM resources or permissions change
4. For a wide or unfamiliar area, launch an `explore` subagent with a concrete question instead of reading file by file.

### Git state

- `git status --short` and `git branch --show-current`. Work on this branch as it is; see Git.
- Record the base branch. The default is `master` (`origin/HEAD`). If the current branch was cut from another feature branch, ask which base the review should diff against.
- Record the uncommitted and untracked files that exist before Step 3 as the baseline. Leave them untouched unless the task requires editing them. Step 5 tells the reviewer which of them are not part of this task.

## Step 2: Clarify

Check each item against the sources. Anything not specified and not derivable from a direct peer is a question.

- Public entry point and names: method on `*PubNub` or on a service, setter names, `WithContext` twin, `PN…` response type names
- File names and unexported builder names (AGENTS.md: the task defines these)
- HTTP method, path, query parameters, headers, body schema, encoding
- Response schema, success status codes, error statuses and how they surface
- Validation: required fields, ranges, mutually exclusive options, and exact error text when it is new
- Defaults and unset-versus-zero semantics (is `0`, `""`, or `nil` a real value or "not set")
- Backward compatibility: any change to an existing exported signature or behavior
- PAM and token requirements; event or subscribe impact
- Test scope: unit tests always; whether e2e tests or example snippets are expected (AGENTS.md: the task defines this)
- Docs: README or `.ai/knowledge` updates when a stated fact or rule changes

How to ask:

- Batch every question into one `AskQuestion` call. Put the recommended option first, based on peers. Ask open-ended questions in text in the same message.
- Do not start Step 3 while a question that changes API shape, wire format, or behavior is open.
- With no open questions, post a short plan and continue without waiting: files to add or change, the public API, the tests, and any doc updates. List assumptions explicitly so the user can correct them.

## Step 3: Implement

### Implementation rules

Guidelines from the repository owner. They override SDK conventions when the two conflict.

- Don't overcomment the functionality. Comments should be related to the work done and be helpful. Don't use comments as a source of information what you did.

<!-- Add further implementation guidelines as bullets above this line. -->

### SDK conventions

- Follow the closest peer for the builder shape: unexported `…Builder` and `…Opts` embedding `endpointOpts`; `new…Builder` delegating to `new…BuilderWithContext`; fluent setters; `GetLogParams`; `Execute` calling `LogUserInput`, then `executeRequest`, then the response parser.
- Exported setters and types get a one-line godoc in the peer's voice. Where it applies, state the contract: required or optional, units, ranges, unset behavior.
- Validation uses `newValidationError(o, Str…)`. Reuse existing `Str…` constants; add new ones next to their siblings in `pubnub.go`. Never change existing message text.
- A new operation needs the iota constant appended in its group, a `String()` case, and a telemetry case. Never reorder existing constants.
- When "unset" differs from the zero value on the wire, track it with a private `set…` flag, as peers do (`setPayload`, `setIfMatchETag`).
- Public API changes are additive. A change that breaks an existing exported signature or behavior needs explicit approval from Step 2.
- New examples and tests use `NewConfigWithUserId`. Acronyms on new exported identifiers stay `UUID`, `HTTP`, `TLS`, `TTL`.
- No hardcoded keys, tokens, or secrets. Use `NewDemoConfig`, `tests/stubs`, or environment variables with a skip.
- New goroutines stop on context cancel and client destroy. Shared maps and listener state stay guarded.
- Keep the diff scoped to the task: no drive-by refactors, renames, or formatting of untouched code.
- Do not bump `Version`, edit the changelog, or edit `.pubnub.yml`. Release tooling owns these.
- Update `README.md` only when the change makes a stated fact or quickstart step false. Update `.ai/knowledge/*` only when the change adds or changes a rule a note states. `go-pr-review` checks both.
- Format every changed `.go` file with `gofmt -w <files>` before Step 4. The repository's gofmt hook runs only when the turn ends, which is too late for the review loop.

## Step 4: Test

### What to cover

Each behavior the change adds or alters gets a unit test in package `pubnub`, using `testify/assert` and `tests/helpers` as peers do:

- Path, query, headers, and body: asserted on `buildPath`, `buildQuery`, `buildHeaders`, `buildBody`
- Every new validation branch, asserting the full error text
- Response parsing: success, empty or partial payload, malformed JSON, and error statuses, through `tests/stubs` or a custom `RoundTripper`
- Unset versus zero-value fields, when the wire format distinguishes them
- Operation `String()` value and telemetry mapping, when peers test them

Add e2e tests or example snippets only when Step 2 put them in scope. Live tests skip when their keys are unset; follow the peer's skip pattern.

### Commands

`make` and `bash` may be missing on Windows. Use the direct equivalents and confirm against `Makefile` and `.github/workflows/` when in doubt.

| Purpose | Command |
|---------|---------|
| Compile code and tests | `go build ./...` and `go test -count=1 -run '^$' ./ ./crypto/... ./utils/ ./tests/helpers/ ./tests/e2e/ ./examples/snippets/api/` |
| Targeted tests | `go test -count=1 -v -run '^TestSetEntity' ./` |
| Root package (`make test-quick`) | `gotestsum --format=testname -- -v ./`, fallback `go test -count=1 ./` |
| Short run (`make test-dev`) | `gotestsum --format=testname -- -timeout=30s -short ./` |
| Other packages | `go test -count=1 ./utils/`, `./tests/helpers/`, `./crypto/...`, `./pnerr/` |
| Vet (CI vets library packages, not `_test.go`) | `go vet <changed library packages>` |
| Format | `gofmt -l <changed .go files>` prints nothing |
| E2E | `go test -count=1 -v ./tests/e2e/` |
| Example snippets | `go test -count=1 ./examples/snippets/api/` |
| Full CI shape | `make test` (needs `make` and `bash`) |

Order: compile, gofmt, vet, targeted tests, root package, then every other affected package.

- Do not run `go test ./...`. It includes `examples/snippets` and `examples/cli`, which CI does not test, and `go test` vet already fails on `examples/snippets/onboarding.go`.
- Live tests (e2e and snippets) need keys. Check presence without printing values: `Test-Path .env` and `[bool]$env:PUBLISH_KEY`, plus `DS_PUBLISH_KEY`, `DS_SUBSCRIBE_KEY`, `DS_SECRET_KEY` for DataSync. Without keys, report them as not run; do not count a skip as a pass.
- CI runs utils, helpers, and e2e with `-race`. When the change touches goroutines or shared state, run `go test -race` on the affected package. If the race detector is unavailable locally (no cgo or gcc), report it as not run.
- On Windows the repository hook replaces Go test output with a `GO TEST SUMMARY`. Treat it as the result. Read `raw_log` only when a failure message is missing.

### When tests fail

1. Decide which side is wrong by checking the requirements list and the peer, not by checking which edit is easier.
   - The test asserts the requirement correctly: fix the implementation.
   - The test asserts something the requirement does not say: fix the test.
2. Never weaken an assertion, delete a test, or add a skip to get a pass.
3. Change an existing test's expectation only when the task deliberately changes that behavior, and name it in the summary.
4. Failure unrelated to the diff: confirm it touches nothing changed, leave it, and report it as pre-existing.
5. Same failure after three fix attempts, or fixing it needs a requirement change: stop and ask.

Step 4 is complete when every test in scope passes or is reported with a reason.

## Step 5: Review loop

Each loop is one fresh `go-pr-review` run plus triage. Run at most 5 loops.

### Launch

Gather these before every loop:

- Absolute repository path, current branch, and the base from Step 1
- `git status --porcelain` and `git diff --name-only <base>`. List untracked new files explicitly, because `go-pr-review` ignores untracked files unless told otherwise. Do not stage them to make them visible.
- Baseline files from Step 1 that this task did not touch
- The requirements list and the decisions from Step 2

Launch through `Task` with `subagent_type: "generalPurpose"`, `run_in_background: false`, and description `go-pr-review loop <n>`. Never use `resume`: every loop needs a fresh context. Do not send the implementation reasoning or earlier findings, so each review stays independent.

Prompt template:

```text
Full repository path: <absolute path>
Read <absolute path>/.cursor/skills/go-pr-review/SKILL.md and follow it exactly.

Target: the current branch `<branch>`. Base: `<base>`.
The change is partly or fully uncommitted. Include uncommitted changes.
These untracked files are part of the change and must be reviewed:
- <path>
These uncommitted files existed before this task and are not part of the change; ignore them:
- <path, or "None">

The change implements <JIRA-KEY>: <title>. Judge correctness against these requirements:
<numbered requirements and agreed decisions>

Read-only: do not edit files. Run only Git commands that read (status, diff, log, show). Never run a Git command that changes state: no checkout, switch, branch creation, add, restore, reset, stash, commit, fetch, pull, or push.
Return only the review, in the skill's Output format.
```

### Triage

Verify every Blocker, Other issue, Tests gap, and Docs gap against the code, the requirements, and the peers before acting. The reviewer can be wrong.

| Verdict | Action |
|---------|--------|
| Valid | Fix it. |
| Invalid | Record the evidence: file and line, requirement, or peer that contradicts it. |
| Pre-existing or out of scope | Record it, do not fix it, and list it in the summary. |
| Needs a user decision (API shape, behavior, scope) | Pause the loop and ask. |

If a later loop raises a finding already rejected and brings no new evidence, keep the rejection.

After fixes: gofmt, the targeted tests, and the root package again, then start the next loop.

### Stop

- Stop when the review answers Merge: Yes and every remaining finding is Invalid or Pre-existing.
- Stop after loop 5. Fix any Valid findings it raised and re-run the tests, then mark those fixes as not re-reviewed in the summary.

Keep a ledger per loop: the verdict, findings by category, what was fixed, and what was rejected and why.

## Step 6: Summary

Reply in chat only; do not write summary files. Lead with the outcome. Spell out requirements in words instead of referring to their numbers.

```markdown
## <JIRA-KEY>: <ticket title>

<One or two sentences: what now works, and the final review verdict.>

### What changed
- Public API and behavior, with the main identifiers in backticks

### Files
- `path` — role in the change

### Decisions and assumptions
- Answers from Step 2 and any assumption the user did not confirm

### Tests
- Added tests and what each group covers
- `<command>` — passed / failed / not run (reason)

### Review loops: <n> of 5
| Loop | Verdict | Blockers | Other issues | Fixed | Rejected |
|------|---------|----------|--------------|-------|----------|

- Rejected findings, one line each with the reason
- Fixes after the last loop that were not re-reviewed

### Open items
- Unresolved findings, pre-existing failures, live tests not run, follow-ups such as snippets, docs, or server questions

### Next steps
- The changes are uncommitted on the current branch `<branch>`. Staging, committing, and pushing are left to the user.
- Offer the `pr-description` skill for the PR title and body.
```

## Stop and ask

Ask the user rather than guess when:

- The requirements and the server behavior or a peer contradict each other
- Fulfilling the task needs a breaking change that Step 2 did not approve
- A review finding needs a product or API decision
- Credentials, the ticket, or the documentation are unreachable
- The work grows beyond the ticket's scope
