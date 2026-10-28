#!/usr/bin/env python3
"""gofmt Go files once, after an agent turn, instead of after every edit.

Cursor fires ``afterFileEdit`` once per file. It has no batch-edit hook.
Sleeping inside ``afterFileEdit`` only delays the next edit; it cannot tell
when the batch is finished. This script records eligible ``.go`` paths while
the agent edits, then runs ``gofmt -w`` a single time from the ``stop`` hook
when that turn ends.

``gofmt`` rewrites each file in full. A syntax error leaves that file
unchanged. The agent does not see the formatted bytes during the same turn,
because formatting happens after the turn's edits are done.

Usage:
  python .cursor/hooks/gofmt_batch.py --record     # afterFileEdit
  python .cursor/hooks/gofmt_batch.py --format     # stop
  python .cursor/hooks/gofmt_batch.py --self-test
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
import time
from pathlib import Path
from typing import Any


MAX_PATHS_PER_GOFMT = 40
STALE_SECONDS = 24 * 3600


def emit(payload: dict[str, Any]) -> None:
    sys.stdout.write(json.dumps(payload, ensure_ascii=False) + "\n")
    sys.stdout.flush()


def finish() -> int:
    """afterFileEdit and stop accept an empty object. Do not send a follow-up."""
    emit({})
    return 0


def state_root() -> Path:
    override = os.environ.get("CURSOR_GOFMT_BATCH_DIR")
    if override:
        root = Path(override)
    else:
        base = (
            os.environ.get("LOCALAPPDATA")
            or os.environ.get("XDG_CACHE_HOME")
            or tempfile.gettempdir()
        )
        root = Path(base) / "cursor-gofmt-batch"
    root.mkdir(parents=True, exist_ok=True)
    return root


def conversation_dir(conversation_id: str) -> Path:
    safe = re.sub(r"[^A-Za-z0-9_-]", "_", conversation_id).strip("_")[:80]
    return state_root() / (safe or "unknown")


def sweep_stale(active: Path | None = None) -> None:
    cutoff = time.time() - STALE_SECONDS
    for child in state_root().iterdir():
        if active is not None and child == active:
            continue
        try:
            if child.is_dir() and child.stat().st_mtime < cutoff:
                shutil.rmtree(child, ignore_errors=True)
        except OSError:
            continue


def load_payload() -> dict[str, Any]:
    raw = sys.stdin.read()
    if not raw.strip():
        return {}
    payload = json.loads(raw)
    return payload if isinstance(payload, dict) else {}


def workspace_roots(payload: dict[str, Any]) -> list[Path]:
    roots = payload.get("workspace_roots") or []
    if isinstance(roots, str):
        roots = [roots]
    resolved: list[Path] = []
    for root in roots:
        if not isinstance(root, str) or not root.strip():
            continue
        try:
            resolved.append(Path(root).resolve())
        except OSError:
            continue
    return resolved


def inside_workspace(path: Path, roots: list[Path]) -> bool:
    if not roots:
        return False
    for root in roots:
        try:
            path.relative_to(root.resolve())
            return True
        except (ValueError, OSError):
            continue
    return False


def eligible(path_text: str, roots: list[Path], *, recorded: bool = False) -> Path | None:
    """Return the resolved path when this edit should be formatted later.

    ``recorded`` is set for paths already accepted during this conversation.
    ``stop`` sometimes arrives without ``workspace_roots``; those paths stay
    eligible because they passed the workspace check when they were recorded.
    """
    if not path_text or not path_text.lower().endswith(".go"):
        return None
    try:
        path = Path(path_text).resolve()
    except OSError:
        return None
    if not path.is_file():
        return None
    lowered_parts = {part.lower() for part in path.parts}
    if "vendor" in lowered_parts or ".cursor" in lowered_parts:
        return None
    if roots:
        if not inside_workspace(path, roots):
            return None
    elif not recorded:
        return None
    return path


def record_path(conversation_id: str, path: Path) -> None:
    directory = conversation_dir(conversation_id)
    directory.mkdir(parents=True, exist_ok=True)
    name = hashlib.sha256(str(path).encode("utf-8")).hexdigest()[:16]
    (directory / name).write_text(str(path), encoding="utf-8")


def take_pending(conversation_id: str) -> list[str]:
    directory = conversation_dir(conversation_id)
    if not directory.is_dir():
        return []
    paths: list[str] = []
    for entry in list(directory.iterdir()):
        if not entry.is_file():
            continue
        try:
            text = entry.read_text(encoding="utf-8").strip()
        except OSError:
            text = ""
        try:
            entry.unlink()
        except OSError:
            pass
        if text:
            paths.append(text)
    try:
        directory.rmdir()
    except OSError:
        pass
    # Preserve first-seen order while dropping duplicates.
    seen: set[str] = set()
    unique: list[str] = []
    for path in paths:
        if path in seen:
            continue
        seen.add(path)
        unique.append(path)
    return unique


def run_gofmt(paths: list[Path]) -> str:
    gofmt = shutil.which("gofmt")
    if gofmt is None:
        return "gofmt is not on PATH; files were left unchanged."
    errors: list[str] = []
    for start in range(0, len(paths), MAX_PATHS_PER_GOFMT):
        chunk = [str(path) for path in paths[start : start + MAX_PATHS_PER_GOFMT]]
        completed = subprocess.run(
            [gofmt, "-w", *chunk],
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            check=False,
        )
        if completed.returncode != 0:
            detail = (completed.stderr or completed.stdout or "").strip()
            errors.append(detail or f"gofmt exited {completed.returncode}")
    return "\n".join(errors)


def record_main(payload: dict[str, Any]) -> int:
    conversation_id = str(payload.get("conversation_id") or "unknown")
    sweep_stale(conversation_dir(conversation_id))
    path = eligible(str(payload.get("file_path") or ""), workspace_roots(payload))
    if path is not None:
        record_path(conversation_id, path)
    return finish()


def format_main(payload: dict[str, Any]) -> int:
    conversation_id = str(payload.get("conversation_id") or "unknown")
    roots = workspace_roots(payload)
    pending = [
        eligible(path, roots, recorded=True) for path in take_pending(conversation_id)
    ]
    paths = [path for path in pending if path is not None]
    if paths:
        errors = run_gofmt(paths)
        if errors:
            print(f"GOFMT BATCH: {errors}", file=sys.stderr)
    return finish()


def hook_main(mode: str) -> int:
    try:
        payload = load_payload()
    except json.JSONDecodeError:
        return finish()
    event = str(payload.get("hook_event_name") or "")
    if mode == "record" or event == "afterFileEdit":
        return record_main(payload)
    if mode == "format" or event == "stop":
        return format_main(payload)
    return finish()


def self_test() -> int:
    failed = 0
    work = Path(tempfile.mkdtemp(prefix="gofmt-batch-"))
    state = work / "state"
    os.environ["CURSOR_GOFMT_BATCH_DIR"] = str(state)
    try:
        outside = work / "outside.go"
        outside.write_text("package p\n", encoding="utf-8")
        repo = work / "repo"
        repo.mkdir()
        sample = repo / "sample.go"
        sample.write_text("package p\n\nfunc  F( ) int {\nreturn  1\n}\n", encoding="utf-8")
        vendor = repo / "vendor" / "x.go"
        vendor.parent.mkdir()
        vendor.write_text("package x\n\nfunc  V( ) int { return 1 }\n", encoding="utf-8")
        roots = [repo]

        if eligible(str(sample), roots) != sample.resolve():
            print("FAIL expected sample.go to be eligible")
            failed += 1
        if eligible(str(outside), roots) is not None:
            print("FAIL path outside the workspace was eligible")
            failed += 1
        if eligible(str(vendor), roots) is not None:
            print("FAIL vendor path was eligible")
            failed += 1
        if eligible(str(repo / "readme.md"), roots) is not None:
            print("FAIL non-go path was eligible")
            failed += 1
        if eligible(str(sample), []) is not None:
            print("FAIL new path with no workspace was eligible")
            failed += 1
        if eligible(str(sample), [], recorded=True) != sample.resolve():
            print("FAIL recorded path was dropped when stop omitted workspace roots")
            failed += 1

        conversation = "self-test"
        record_path(conversation, sample.resolve())
        record_path(conversation, sample.resolve())
        pending = take_pending(conversation)
        if pending != [str(sample.resolve())]:
            print(f"FAIL pending paths: {pending}")
            failed += 1
        if conversation_dir(conversation).exists():
            print("FAIL pending directory was not cleared")
            failed += 1

        record_path(conversation, sample.resolve())
        formatted = [path for path in (eligible(item, roots) for item in take_pending(conversation)) if path]
        errors = run_gofmt(formatted)
        if errors:
            print(f"FAIL gofmt: {errors}")
            failed += 1
        formatted_text = sample.read_text(encoding="utf-8")
        if formatted_text != "package p\n\nfunc F() int {\n\treturn 1\n}\n":
            print("FAIL formatted source")
            print(formatted_text)
            failed += 1

        broken = repo / "broken.go"
        broken.write_text("package p\nfunc\n", encoding="utf-8")
        original = broken.read_text(encoding="utf-8")
        broken_errors = run_gofmt([broken])
        if not broken_errors:
            print("FAIL expected gofmt to reject a syntax error")
            failed += 1
        if broken.read_text(encoding="utf-8") != original:
            print("FAIL gofmt rewrote a file it could not parse")
            failed += 1
    finally:
        shutil.rmtree(work, ignore_errors=True)
        os.environ.pop("CURSOR_GOFMT_BATCH_DIR", None)

    if failed:
        print(f"self-test failed: {failed} check(s)")
        return 1
    print("self-test passed")
    return 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Batch gofmt for Cursor agent edits.")
    parser.add_argument("--record", action="store_true", help="Record afterFileEdit paths.")
    parser.add_argument("--format", action="store_true", help="gofmt paths recorded for this turn.")
    parser.add_argument("--self-test", action="store_true")
    args = parser.parse_args(argv)

    if args.self_test:
        return self_test()
    if args.record and args.format:
        print("GOFMT BATCH: pass only one of --record or --format", file=sys.stderr)
        return finish()
    mode = "format" if args.format else "record" if args.record else ""
    return hook_main(mode)


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # noqa: BLE001 - hook must fail open
        if "--self-test" in sys.argv:
            raise
        print(f"GOFMT BATCH: {exc}", file=sys.stderr)
        emit({})
        raise SystemExit(0)
