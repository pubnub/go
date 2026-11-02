#!/usr/bin/env python3
"""Cursor preToolUse hook: wrap Go test runs and keep only failures.

Windows only. Cursor loads this project hook on every platform, including
macOS, Linux, and cloud agents. The ``_windows`` name does not filter by OS.
On non-Windows hosts ``--hook`` reads stdin and returns ``{"permission":"allow"}``
so the original command runs unchanged. Wrapping uses ``powershell.exe`` and
``%LOCALAPPDATA%``, which exist on Windows.

preToolUse cannot replace a Shell tool's output after the fact. The command is
rewritten so the process the agent runs prints failed tests, their errors, and
the final result. The full log is still written to disk.

Matches plain ``go test``, ``gotestsum``, Makefile test targets, and
``run-tests-gotestsum.sh`` (same entry points as CI / this repo's Makefile).

Usage:
  python .cursor/hooks/go_test_wrap_windows.py --hook
  python .cursor/hooks/go_test_wrap_windows.py --run --cmd-file <path>
  python .cursor/hooks/go_test_wrap_windows.py --self-test
"""

from __future__ import annotations

import argparse
import json
import os
import re
import subprocess
import sys
import tempfile
import time
import uuid
from collections import deque
from pathlib import Path
from typing import Any


SCRIPT_PATH = Path(__file__).resolve()
MAX_FAILED = 40
MAX_ERRORS_PER_TEST = 10
MAX_ORPHAN_ERRORS = 12
MAX_LINE_CHARS = 400
MAX_TAIL_LINES = 12
MAX_LOG_BYTES = 32 * 1024 * 1024
START_BANNER = (
    "GO TEST WRAP: running Go tests; model will receive failures and the final "
    "result, not the raw log."
)

OTHER_WRAP_MARKERS = (
    "compile_wrap.py",
    "ue_test_wrap.py",
    "go_test_wrap_windows.py",
    "go_test_wrap.py",
    "CURSOR_COMPILE_WRAP",
    "CURSOR_UE_TEST_WRAP",
    "CURSOR_GO_TEST_WRAP",
)

VIEWER_PREFIX = re.compile(
    r"(?i)^\s*(?:&\s*)?(?:"
    r"Get-Content|Get-Item|Get-ChildItem|Get-Command|Get-Help|"
    r"Select-String|Out-File|Set-Content|"
    r"gc\b|cat\b|type\b|more\b|less\b|tail\b|head\b|"
    r"findstr\b|find\b|rg\b|grep\b|where(?:\.exe)?|"
    r"notepad(?:\+\+)?|code\b|cursor\b|vim\b|emacs\b|"
    r"dir\b|ls\b|echo\b|Write-Output|Write-Host"
    r")\b"
)

GO_TEST_PATTERNS = (
    re.compile(r"(?i)(?:^|[\s;&|])(?:go\s+test|gotestsum)\b"),
    re.compile(r"(?i)run-tests-gotestsum\.sh"),
    re.compile(
        r"(?i)\bmake\b(?:\s+[^\n|;&]*)?\btest(?:-quick|-dev|-utils|-helpers|-integration|-coverage|-format)?\b"
    ),
)

GO_NON_TEST = re.compile(
    r"(?i)(?:^|[\s;&|])go\s+(?:build|mod|install|get|run|fmt|vet|generate|doc|list|clean|work|env|version|tool)\b"
)

RUN_LINE = re.compile(r"^=== RUN\s+(?P<name>\S+)")
FAIL_HEADER = re.compile(r"^--- FAIL:\s+(?P<name>\S+)\s+\(")
PASS_HEADER = re.compile(r"^--- PASS:\s+(?P<name>\S+)\s+\(")
SKIP_HEADER = re.compile(r"^--- SKIP:\s+(?P<name>\S+)\s+\(")
PKG_FAIL = re.compile(r"^FAIL\s+(?P<pkg>\S+)")
PKG_OK = re.compile(r"^ok\s+(?P<pkg>\S+)")
FAIL_INDENT = re.compile(r"^\s+(?P<file>[^\s:]+):(?P<line>\d+):\s+(?P<msg>.*)$")
FAIL_INDENT_TAB = re.compile(r"^\t(?P<msg>.*)$")
PANIC_LINE = re.compile(r"^panic:\s+(?P<msg>.*)$")
GOTESTSUM_CASE_FAIL = re.compile(
    r"^[✖×xX✗]\s+(?:(?P<pkg>[\w./-]+)\.)?(?P<name>\S+)\s+\("
)
GOTESTSUM_CASE_PASS = re.compile(
    r"^[✓✔√]\s+(?:(?P<pkg>[\w./-]+)\.)?(?P<name>\S+)\s+\("
)
GOTESTSUM_DONE = re.compile(
    r"(?i)^DONE\s+(?P<count>\d+)\s+tests?(?:,\s*(?P<failures>\d+)\s+failures?)?"
)
BUILD_FAIL = re.compile(r"(?i)^FAIL\s*$")
NOISE = re.compile(r"(?i)(?:no tests to run|testing: warning|coverage:)")


def emit(payload: dict[str, Any]) -> None:
    sys.stdout.write(json.dumps(payload, ensure_ascii=False) + "\n")
    sys.stdout.flush()


def allow() -> int:
    emit({"permission": "allow"})
    return 0


def is_windows() -> bool:
    """True only for native Windows. Cursor still invokes this hook elsewhere."""
    return sys.platform == "win32"


def wrap_state_dir() -> Path:
    root = Path(os.environ.get("LOCALAPPDATA") or tempfile.gettempdir()) / "cursor-go-test-wrap"
    root.mkdir(parents=True, exist_ok=True)
    (root / "logs").mkdir(exist_ok=True)
    cutoff = time.time() - 24 * 3600
    for stale in root.glob("cmd-*.json"):
        try:
            if stale.stat().st_mtime < cutoff:
                stale.unlink()
        except OSError:
            pass
    return root


def is_already_wrapped(command: str) -> bool:
    return any(marker in command for marker in OTHER_WRAP_MARKERS)


def is_go_test_command(command: str) -> bool:
    if not command or not command.strip():
        return False
    if is_already_wrapped(command):
        return False
    if VIEWER_PREFIX.match(command):
        return False
    if not any(pattern.search(command) for pattern in GO_TEST_PATTERNS):
        return False
    if GO_NON_TEST.search(command) and not re.search(
        r"(?i)(?:^|[\s;&|])(?:go\s+test|gotestsum)\b", command
    ):
        return False
    return True


def parse_tool_input(payload: dict[str, Any]) -> dict[str, Any]:
    tool_input = payload.get("tool_input") or {}
    if isinstance(tool_input, str):
        try:
            parsed = json.loads(tool_input)
            tool_input = parsed if isinstance(parsed, dict) else {"command": tool_input}
        except json.JSONDecodeError:
            tool_input = {"command": tool_input}
    if not isinstance(tool_input, dict):
        return {}
    return dict(tool_input)


def command_from_input(tool_input: dict[str, Any], payload: dict[str, Any]) -> str:
    for key in ("command", "cmd"):
        value = tool_input.get(key)
        if isinstance(value, str) and value.strip():
            return value
    value = payload.get("command")
    return value if isinstance(value, str) else ""


def cwd_from_input(tool_input: dict[str, Any], payload: dict[str, Any]) -> str:
    for source in (tool_input, payload):
        for key in ("working_directory", "workingDirectory", "cwd"):
            value = source.get(key)
            if isinstance(value, str) and value.strip():
                return value
    return ""


def truncate(text: str, limit: int = MAX_LINE_CHARS) -> str:
    text = " ".join(text.split())
    if len(text) <= limit:
        return text
    return text[: limit - 3] + "..."


class GoTestRun:
    def __init__(self) -> None:
        self.passed = 0
        self.failed = 0
        self.skipped = 0
        self.failed_tests: list[dict[str, Any]] = []
        self.orphan_errors: list[str] = []
        self.current: dict[str, Any] | None = None
        self.in_failure_block = False
        self.failed_packages: list[str] = []
        self.done_total: int | None = None
        self.done_failures: int | None = None
        self.saw_build_fail = False
        self.tail: deque[str] = deque(maxlen=MAX_TAIL_LINES)

    def _remember_error(self, text: str) -> None:
        text = truncate(text)
        if not text or NOISE.search(text):
            return
        if self.current is not None:
            errors: list[str] = self.current["errors"]
            if text not in errors and len(errors) < MAX_ERRORS_PER_TEST:
                errors.append(text)
            return
        if text not in self.orphan_errors and len(self.orphan_errors) < MAX_ORPHAN_ERRORS:
            self.orphan_errors.append(text)

    def _close_current(self) -> None:
        self.in_failure_block = False
        current = self.current
        self.current = None
        if current is None:
            return
        name = current["name"] or "(unnamed test)"
        entry = {
            "name": name,
            "pkg": current.get("pkg") or "",
            "errors": list(current["errors"]),
        }
        self.failed += 1
        self.failed_tests.append(entry)

    def _start_failed(self, name: str, pkg: str = "") -> None:
        self._close_current()
        self.current = {"name": name, "pkg": pkg, "errors": []}
        self.in_failure_block = True

    def consume(self, line: str) -> None:
        stripped = line.rstrip("\n").rstrip()
        if not stripped:
            return

        if RUN_LINE.match(stripped):
            return

        fail_header = FAIL_HEADER.match(stripped)
        if fail_header:
            self._start_failed(fail_header.group("name"))
            return

        if PASS_HEADER.match(stripped) or SKIP_HEADER.match(stripped):
            self._close_current()
            if SKIP_HEADER.match(stripped):
                self.skipped += 1
            else:
                self.passed += 1
            return

        gs_fail = GOTESTSUM_CASE_FAIL.match(stripped)
        if gs_fail:
            pkg = gs_fail.group("pkg") or ""
            name = gs_fail.group("name")
            label = f"{pkg}.{name}" if pkg else name
            self._start_failed(label, pkg)
            return

        gs_pass = GOTESTSUM_CASE_PASS.match(stripped)
        if gs_pass:
            self._close_current()
            self.passed += 1
            return

        done = GOTESTSUM_DONE.match(stripped)
        if done:
            self.done_total = int(done.group("count"))
            failures = done.group("failures")
            if failures is not None:
                self.done_failures = int(failures)
            return

        pkg_fail = PKG_FAIL.match(stripped)
        if pkg_fail:
            pkg = pkg_fail.group("pkg")
            if pkg not in self.failed_packages:
                self.failed_packages.append(pkg)
            return

        if PKG_OK.match(stripped):
            return

        if BUILD_FAIL.match(stripped):
            self.saw_build_fail = True
            return

        panic = PANIC_LINE.match(stripped)
        if panic:
            self.saw_build_fail = True
            self._remember_error(f"panic: {panic.group('msg')}")
            return

        if self.in_failure_block and self.current is not None:
            if stripped.startswith("--- "):
                self._close_current()
                return
            indent = FAIL_INDENT.match(stripped) or FAIL_INDENT_TAB.match(stripped)
            if indent:
                msg = indent.group("msg") if "msg" in indent.groupdict() else stripped.strip()
                if "file" in indent.groupdict() and indent.group("file"):
                    msg = f"{indent.group('file')}:{indent.group('line')}: {msg}"
                self._remember_error(msg)
                return
            if stripped.startswith("\t") or stripped.startswith("    "):
                self._remember_error(stripped.strip())
                return

        if FAIL_HEADER.search(stripped) or "Error Trace" in stripped:
            self._remember_error(stripped)

        if not stripped.startswith("["):
            self.tail.append(truncate(stripped))


def format_summary(
    *,
    command: str,
    exit_code: int,
    duration_sec: float,
    run: GoTestRun,
    log_path: Path,
    raw_lines: int,
) -> str:
    run._close_current()
    failed = (
        exit_code != 0
        or run.failed > 0
        or run.saw_build_fail
        or bool(run.failed_packages)
        or (run.done_failures or 0) > 0
    )
    status = "FAILED" if failed else "SUCCEEDED"
    shown = run.failed_tests[:MAX_FAILED]
    omitted = max(0, len(run.failed_tests) - len(shown))

    lines = [
        "=== GO TEST SUMMARY (raw test log stripped to save context) ===",
        f"status: {status}",
        f"exit_code: {exit_code}",
        f"duration_sec: {duration_sec:.1f}",
        f"tests_passed: {run.passed}",
        f"tests_failed: {run.failed}",
        f"tests_skipped: {run.skipped}",
    ]
    if run.done_total is not None:
        lines.append(f"gotestsum_done_total: {run.done_total}")
    if run.done_failures is not None:
        lines.append(f"gotestsum_done_failures: {run.done_failures}")
    if run.failed_packages:
        lines.append(f"failed_packages: {', '.join(run.failed_packages[:20])}")
    lines.extend(
        [
            f"raw_log_lines: {raw_lines}",
            f"raw_log: {log_path}",
            f"original_command: {truncate(command, 300)}",
            "",
            "Use this summary as the Go test result. Do not re-run tests just to inspect failures.",
            "Open the cited source files. Read raw_log only if a failure message is missing here.",
        ]
    )

    if shown:
        lines.append("")
        lines.append("FAILED TESTS:")
        for index, test in enumerate(shown, start=1):
            label = test["name"]
            if test.get("pkg") and test["pkg"] not in label:
                label = f"{test['pkg']} {label}"
            lines.append(f"{index}. {truncate(label)}")
            for error in test["errors"]:
                lines.append(f"   - {error}")
            if not test["errors"]:
                lines.append("   - (no parsed error text; see raw_log)")
        if omitted:
            lines.append(f"... {omitted} more failed tests omitted; see raw_log")

    if run.orphan_errors and (failed or not shown):
        lines.append("")
        lines.append("ERRORS:")
        for error in run.orphan_errors:
            lines.append(f"- {error}")

    tail = [item for item in run.tail if item.strip()]
    if failed and not shown and not run.orphan_errors and tail:
        lines.append("")
        lines.append("TAIL (no structured test failures parsed; last output):")
        lines.extend(f"| {item}" for item in tail)

    if not failed:
        lines.append("")
        lines.append("tests reported success")

    lines.append("=== END GO TEST SUMMARY ===")
    return "\n".join(lines) + "\n"


def run_tests(command: str, cwd: str | None, log_path: Path) -> tuple[int, str]:
    run = GoTestRun()
    raw_lines = 0
    log_bytes = 0

    popen_kwargs: dict[str, Any] = {
        "args": [
            "powershell.exe",
            "-NoProfile",
            "-NonInteractive",
            "-ExecutionPolicy",
            "Bypass",
            "-Command",
            command,
        ],
        "stdout": subprocess.PIPE,
        "stderr": subprocess.STDOUT,
        "cwd": cwd or None,
        "env": {**os.environ, "CURSOR_GO_TEST_WRAP": "1"},
        "text": True,
        "encoding": "utf-8",
        "errors": "replace",
        "bufsize": 1,
    }
    if os.name == "nt":
        popen_kwargs["creationflags"] = getattr(subprocess, "CREATE_NO_WINDOW", 0)

    started = time.perf_counter()
    proc = subprocess.Popen(**popen_kwargs)
    assert proc.stdout is not None

    print(START_BANNER, flush=True)

    with log_path.open("w", encoding="utf-8", errors="replace") as log_file:
        log_file.write(f"# original_command: {command}\n")
        log_file.write(f"# cwd: {cwd or os.getcwd()}\n\n")
        for line in proc.stdout:
            raw_lines += 1
            encoded = line.encode("utf-8", "replace")
            if log_bytes < MAX_LOG_BYTES:
                writable = encoded[: MAX_LOG_BYTES - log_bytes]
                log_file.write(writable.decode("utf-8", "replace"))
                log_bytes += len(writable)
                if log_bytes >= MAX_LOG_BYTES:
                    log_file.write("\n# truncated: log exceeded 32 MiB\n")
            run.consume(line)

    exit_code = proc.wait()
    duration = time.perf_counter() - started
    summary = format_summary(
        command=command,
        exit_code=exit_code,
        duration_sec=duration,
        run=run,
        log_path=log_path,
        raw_lines=raw_lines,
    )
    return exit_code, summary


def hook_main() -> int:
    try:
        raw = sys.stdin.read()
    except Exception:
        raw = ""
    if not is_windows():
        return allow()
    try:
        payload = json.loads(raw) if raw.strip() else {}
    except json.JSONDecodeError:
        return allow()

    if payload.get("tool_name") and payload.get("tool_name") != "Shell":
        return allow()

    tool_input = parse_tool_input(payload)
    command = command_from_input(tool_input, payload)
    if not is_go_test_command(command):
        return allow()

    cwd = cwd_from_input(tool_input, payload)
    cmd_path = wrap_state_dir() / f"cmd-{uuid.uuid4().hex}.json"
    cmd_path.write_text(
        json.dumps(
            {
                "command": command,
                "cwd": cwd,
                "conversation_id": payload.get("conversation_id"),
            },
            ensure_ascii=False,
            indent=2,
        ),
        encoding="utf-8",
    )

    wrap_command = subprocess.list2cmdline(
        [sys.executable, str(SCRIPT_PATH), "--run", "--cmd-file", str(cmd_path)]
    )
    updated = dict(tool_input)
    updated["command"] = wrap_command
    emit(
        {
            "permission": "allow",
            "updated_input": updated,
            "agent_message": (
                "Go test command wrapped: you will receive failed tests, their errors, "
                "and the final result instead of the raw test log. Full log is saved on disk."
            ),
        }
    )
    return 0


def run_main(cmd_file: Path) -> int:
    if not is_windows():
        print(
            "GO TEST WRAP: go_test_wrap_windows.py runs the wrapped command on Windows only.",
            file=sys.stderr,
        )
        return 1
    try:
        payload = json.loads(cmd_file.read_text(encoding="utf-8-sig"))
    except (OSError, json.JSONDecodeError) as exc:
        print(f"GO TEST WRAP: failed to read command file {cmd_file}: {exc}", file=sys.stderr)
        return 1
    finally:
        try:
            cmd_file.unlink(missing_ok=True)
        except OSError:
            pass

    command = payload.get("command") or ""
    cwd = payload.get("cwd") or None
    if cwd and not Path(cwd).is_dir():
        cwd = None
    if not command.strip():
        print("GO TEST WRAP: empty command", file=sys.stderr)
        return 1

    stamp = time.strftime("%Y%m%d-%H%M%S")
    log_path = wrap_state_dir() / "logs" / f"{stamp}-{uuid.uuid4().hex[:8]}.log"
    exit_code, summary = run_tests(command, cwd, log_path)
    sys.stdout.write(summary)
    sys.stdout.flush()
    return exit_code


def self_test() -> int:
    test_cases = [
        "go test -v ./...",
        "go test -short ./utils/",
        r"gotestsum --format=testname -- -v ./",
        r"bash ./scripts/run-tests-gotestsum.sh -mod=mod testname",
        "make test-quick",
        "make test-integration",
        "make test",
    ]
    non_test_cases = [
        "go build ./...",
        "go mod tidy",
        "go vet ./...",
        "git status",
        "Get-Content go_test_wrap.py",
        "rg \"go test\" Makefile",
        "echo go test ./...",
        f"{sys.executable} {SCRIPT_PATH} --run --cmd-file x.json",
        "make help",
        "make install-tools",
    ]
    failed = 0
    for command in test_cases:
        if not is_go_test_command(command):
            print(f"FAIL expected go test: {command}")
            failed += 1
    for command in non_test_cases:
        if is_go_test_command(command):
            print(f"FAIL expected skip: {command}")
            failed += 1

    sample = [
        "=== RUN   TestPass",
        "--- PASS: TestPass (0.01s)",
        "=== RUN   TestFail",
        "--- FAIL: TestFail (0.02s)",
        "    grant_token_request_test.go:42: expected nil, got error",
        "FAIL",
        "FAIL\tgithub.com/pubnub/go/v7\t0.123s",
        "✓  github.com/pubnub/go/v7.TestOther (0.01s)",
        "✖  github.com/pubnub/go/v7.TestFail (0.02s)",
        "DONE 12 tests, 1 failure in 3.4s",
    ]
    run = GoTestRun()
    for line in sample:
        run.consume(line)
    if run.failed < 1 or run.passed < 1:
        print(f"FAIL counts failed={run.failed} passed={run.passed}")
        failed += 1
    if run.done_total != 12 or run.done_failures != 1:
        print(f"FAIL done markers total={run.done_total} failures={run.done_failures}")
        failed += 1
    if not run.failed_tests or "expected nil" not in " ".join(run.failed_tests[0]["errors"]):
        print(f"FAIL missing error text: {run.failed_tests}")
        failed += 1

    summary = format_summary(
        command=test_cases[0],
        exit_code=1,
        duration_sec=3.4,
        run=run,
        log_path=Path("C:/temp/go-test.log"),
        raw_lines=len(sample),
    )
    if "status: FAILED" not in summary or "TestFail" not in summary:
        print("FAIL summary content")
        print(summary)
        failed += 1

    if failed:
        print(f"self-test failed: {failed} check(s)")
        return 1
    print("self-test passed")
    return 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Wrap Go test runs for Cursor agents.")
    parser.add_argument("--hook", action="store_true", help="Read Cursor preToolUse JSON from stdin.")
    parser.add_argument("--run", action="store_true", help="Execute a previously captured test command.")
    parser.add_argument("--cmd-file", type=Path, help="JSON file written by --hook.")
    parser.add_argument("--self-test", action="store_true")
    args = parser.parse_args(argv)

    if args.self_test:
        return self_test()
    if args.run:
        if not args.cmd_file:
            print("GO TEST WRAP: --cmd-file is required with --run", file=sys.stderr)
            return 1
        return run_main(args.cmd_file)
    return hook_main()


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # noqa: BLE001 - hook must fail open
        if "--run" in sys.argv:
            print(f"GO TEST WRAP: {exc}", file=sys.stderr)
            raise
        emit({"permission": "allow"})
        raise SystemExit(0)
