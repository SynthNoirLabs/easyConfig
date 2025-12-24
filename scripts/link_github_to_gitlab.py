#!/usr/bin/env python3
"""
Post back-links on GitHub issues/PRs to their migrated GitLab counterparts and optionally close them.

This uses the GitHub REST API (via direct HTTPS) and expects:
- GITHUB_TOKEN in env (repo scope)

Inputs:
- One or more mapping JSON files produced by:
  - scripts/migrate_github_issues_to_gitlab.py
  - scripts/migrate_github_prs_to_gitlab.py

Usage:
  export GITHUB_TOKEN=...
  python3 scripts/link_github_to_gitlab.py \
    --repo komod0/easyConfig \
    --map migration-issue-map-open.json \
    --map migration-issue-map-closed.json \
    --map migration-pr-map-all.json \
    --comment --close-open --yes
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import time
import urllib.request
from dataclasses import dataclass
from typing import Any, Dict, Iterable, List, Optional, Tuple


GITHUB_API = "https://api.github.com"


def eprint(*args: object) -> None:
    print(*args, file=sys.stderr)


class HTTPError(RuntimeError):
    pass


def http_json(
    method: str,
    url: str,
    headers: Dict[str, str],
    body: Optional[Dict[str, Any]] = None,
    timeout_s: int = 60,
) -> Any:
    data = None
    if body is not None:
        data = json.dumps(body).encode("utf-8")
        headers = dict(headers)
        headers.setdefault("Content-Type", "application/json")
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=timeout_s) as resp:
            raw = resp.read()
            if not raw:
                return None
            return json.loads(raw.decode("utf-8"))
    except urllib.error.HTTPError as ex:
        raw = ex.read().decode("utf-8", errors="replace")
        raise HTTPError(f"{method} {url} -> {ex.code}\n{raw}") from ex


def gh_headers(token: str) -> Dict[str, str]:
    return {
        "Accept": "application/vnd.github+json",
        "Authorization": f"Bearer {token}",
        "X-GitHub-Api-Version": "2022-11-28",
        "User-Agent": "easyconfig-migrate-script",
    }


@dataclass(frozen=True)
class Target:
    kind: str  # issue|pr
    number: int
    gitlab_url: str
    gitlab_type: str  # issue|mr


def load_targets_from_map(path: str) -> List[Target]:
    data = json.load(open(path, "r", encoding="utf-8"))
    targets: List[Target] = []

    if "issues" in data:
        for entry in data.get("issues", []):
            n = entry.get("github_number")
            url = entry.get("gitlab_url")
            if isinstance(n, int) and isinstance(url, str) and url:
                targets.append(Target(kind="issue", number=n, gitlab_url=url, gitlab_type="issue"))

    if "items" in data:
        for entry in data.get("items", []):
            n = entry.get("github_pr")
            url = entry.get("gitlab_url")
            t = entry.get("type")
            if isinstance(n, int) and isinstance(url, str) and url:
                targets.append(Target(kind="pr", number=n, gitlab_url=url, gitlab_type=("mr" if t == "mr" else "issue")))

    return targets


def get_issue_state(repo: str, number: int, token: str) -> Tuple[str, bool]:
    url = f"{GITHUB_API}/repos/{repo}/issues/{number}"
    issue = http_json("GET", url, gh_headers(token))
    state = str(issue.get("state", "unknown"))
    is_pr = "pull_request" in issue
    return state, is_pr


def post_comment(repo: str, number: int, body: str, token: str, dry_run: bool) -> None:
    if dry_run:
        return
    url = f"{GITHUB_API}/repos/{repo}/issues/{number}/comments"
    http_json("POST", url, gh_headers(token), {"body": body})


def close_issue_or_pr(repo: str, number: int, token: str, dry_run: bool) -> None:
    if dry_run:
        return
    url = f"{GITHUB_API}/repos/{repo}/issues/{number}"
    http_json("PATCH", url, gh_headers(token), {"state": "closed"})


def format_comment(t: Target) -> str:
    if t.kind == "issue":
        return f"Moved to GitLab: {t.gitlab_url}\n\nPlease continue discussion there."
    if t.gitlab_type == "mr":
        return f"Migrated to GitLab Merge Request: {t.gitlab_url}\n\nThis PR will be closed on GitHub; please continue work on GitLab."
    return f"Archived on GitLab: {t.gitlab_url}"


def main() -> int:
    p = argparse.ArgumentParser()
    p.add_argument("--repo", required=True, help="owner/name, e.g. komod0/easyConfig")
    p.add_argument("--map", action="append", default=[], help="Mapping JSON file (repeatable)")
    p.add_argument("--comment", action="store_true")
    p.add_argument("--close-open", action="store_true", help="Close items that are currently open on GitHub")
    p.add_argument("--rate-limit-sleep", type=float, default=0.25)
    p.add_argument("--yes", action="store_true")
    args = p.parse_args()

    token = os.environ.get("GITHUB_TOKEN", "")
    if not token:
        eprint("Missing env var: GITHUB_TOKEN")
        return 2
    if not args.map:
        eprint("No --map provided")
        return 2

    dry_run = not args.yes
    targets: List[Target] = []
    for mp in args.map:
        try:
            targets.extend(load_targets_from_map(mp))
        except FileNotFoundError:
            eprint(f"Skipping missing map: {mp}")
            continue

    # de-dupe by (kind, number)
    uniq: Dict[Tuple[str, int], Target] = {}
    for t in targets:
        uniq[(t.kind, t.number)] = t
    targets = sorted(uniq.values(), key=lambda t: (t.kind, t.number))

    eprint(f"Targets: {len(targets)} (dry_run={dry_run})")

    commented = 0
    closed = 0
    for t in targets:
        state, is_pr = get_issue_state(args.repo, t.number, token)
        # sanity: kind mismatch isn't fatal
        comment_body = format_comment(t)

        if args.comment:
            eprint(f"Commenting on #{t.number} ({'PR' if is_pr else 'issue'}; state={state}) -> {t.gitlab_url}")
            post_comment(args.repo, t.number, comment_body, token, dry_run)
            commented += 1

        if args.close_open and state == "open":
            eprint(f"Closing #{t.number} on GitHub")
            close_issue_or_pr(args.repo, t.number, token, dry_run)
            closed += 1

        time.sleep(args.rate_limit_sleep)

    eprint(f"Done. commented={commented} closed={closed}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

