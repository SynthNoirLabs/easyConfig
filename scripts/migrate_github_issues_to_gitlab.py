#!/usr/bin/env python3
"""
Migrate GitHub issues -> GitLab issues (optionally comments).

Limitations:
- Pull requests are NOT migrated (GitHub PRs are separate objects in GitLab as MRs).
- Assignees/milestones/users may not map 1:1 across platforms.

Requires:
- GITHUB_TOKEN: GitHub token with permission to read issues.
- GITLAB_TOKEN: GitLab token with `api` scope.

Usage:
  python3 scripts/migrate_github_issues_to_gitlab.py \
    --github-repo komod0/easyConfig \
    --gitlab-project komod0/easyconfig \
    --state open \
    --with-comments \
    --close-when-closed \
    --yes
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import time
import urllib.parse
import urllib.request
from dataclasses import dataclass
from typing import Any, Dict, Iterable, List, Optional, Tuple


GITHUB_API = "https://api.github.com"
GITLAB_API = "https://gitlab.com/api/v4"


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


def gl_headers(token: str) -> Dict[str, str]:
    return {
        "PRIVATE-TOKEN": token,
        "User-Agent": "easyconfig-migrate-script",
    }


def paginate_github(url: str, token: str) -> Iterable[Any]:
    page = 1
    while True:
        sep = "&" if "?" in url else "?"
        page_url = f"{url}{sep}per_page=100&page={page}"
        items = http_json("GET", page_url, gh_headers(token))
        if not isinstance(items, list):
            raise RuntimeError(f"Expected list from GitHub API, got: {type(items)}")
        if not items:
            break
        yield from items
        page += 1


@dataclass(frozen=True)
class Issue:
    number: int
    title: str
    body: str
    state: str
    html_url: str
    user_login: str
    created_at: str
    labels: List[str]


def fetch_github_issues(repo: str, token: str, state: str) -> List[Issue]:
    url = f"{GITHUB_API}/repos/{repo}/issues?state={urllib.parse.quote(state)}"
    issues: List[Issue] = []
    for item in paginate_github(url, token):
        # GitHub "issues" API returns PRs too; skip them.
        if "pull_request" in item:
            continue
        issues.append(
            Issue(
                number=int(item["number"]),
                title=str(item["title"]),
                body=str(item.get("body") or ""),
                state=str(item["state"]),
                html_url=str(item["html_url"]),
                user_login=str(item["user"]["login"]),
                created_at=str(item["created_at"]),
                labels=[str(l["name"]) for l in (item.get("labels") or [])],
            )
        )
    return issues


def fetch_github_issue_comments(repo: str, issue_number: int, token: str) -> List[Dict[str, Any]]:
    url = f"{GITHUB_API}/repos/{repo}/issues/{issue_number}/comments"
    return list(paginate_github(url, token))


def gitlab_project_id(project_path: str, token: str) -> int:
    encoded = urllib.parse.quote(project_path, safe="")
    url = f"{GITLAB_API}/projects/{encoded}"
    proj = http_json("GET", url, gl_headers(token))
    return int(proj["id"])


def ensure_gitlab_label(project_id: int, label: str, token: str, dry_run: bool) -> None:
    # GitLab label create is idempotent-ish: if exists it errors. We tolerate that.
    if dry_run:
        return
    url = f"{GITLAB_API}/projects/{project_id}/labels"
    try:
        http_json("POST", url, gl_headers(token), {"name": label, "color": "#428BCA"})
    except HTTPError as ex:
        # "name has already been taken" or "Label already exists" etc.
        if "already been taken" in str(ex) or "Label" in str(ex):
            return
        raise


def create_gitlab_issue(
    project_id: int,
    title: str,
    description: str,
    labels: List[str],
    token: str,
    dry_run: bool,
) -> Optional[Dict[str, Any]]:
    if dry_run:
        return None
    url = f"{GITLAB_API}/projects/{project_id}/issues"
    payload: Dict[str, Any] = {"title": title, "description": description}
    if labels:
        payload["labels"] = ",".join(labels)
    return http_json("POST", url, gl_headers(token), payload)


def create_gitlab_note(project_id: int, issue_iid: int, body: str, token: str, dry_run: bool) -> None:
    if dry_run:
        return
    url = f"{GITLAB_API}/projects/{project_id}/issues/{issue_iid}/notes"
    http_json("POST", url, gl_headers(token), {"body": body})

def close_gitlab_issue(project_id: int, issue_iid: int, token: str, dry_run: bool) -> None:
    if dry_run:
        return
    url = f"{GITLAB_API}/projects/{project_id}/issues/{issue_iid}"
    http_json("PUT", url, gl_headers(token), {"state_event": "close"})


def format_issue_description(issue: Issue) -> str:
    header = (
        f"Migrated from GitHub: {issue.html_url}\n"
        f"Original author: @{issue.user_login}\n"
        f"Original created at: {issue.created_at}\n"
        f"\n---\n\n"
    )
    return header + (issue.body or "")


def format_comment_note(c: Dict[str, Any]) -> str:
    return (
        f"Migrated comment from GitHub: {c.get('html_url','')}\n"
        f"Author: @{c.get('user',{}).get('login','unknown')}\n"
        f"Created at: {c.get('created_at','')}\n\n"
        f"{c.get('body') or ''}"
    )


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--github-repo", required=True, help="owner/name, e.g. komod0/easyConfig")
    parser.add_argument("--gitlab-project", required=True, help="namespace/path, e.g. komod0/easyconfig")
    parser.add_argument("--state", default="open", choices=["open", "closed", "all"])
    parser.add_argument("--with-comments", action="store_true")
    parser.add_argument("--close-when-closed", action="store_true", help="Close GitLab issue if GitHub issue is closed")
    parser.add_argument("--rate-limit-sleep", type=float, default=0.2)
    parser.add_argument("--yes", action="store_true", help="Actually create issues (otherwise dry-run)")
    parser.add_argument("--out", default="migration-issue-map.json", help="Write mapping JSON to this file")
    parser.add_argument(
        "--skip-map",
        action="append",
        default=[],
        help="Path to a previous mapping JSON file; GitHub issue numbers in it will be skipped",
    )

    args = parser.parse_args()

    gh_token = os.environ.get("GITHUB_TOKEN", "")
    gl_token = os.environ.get("GITLAB_TOKEN", "")
    if not gh_token:
        eprint("Missing env var: GITHUB_TOKEN")
        return 2
    if not gl_token:
        eprint("Missing env var: GITLAB_TOKEN")
        return 2

    dry_run = not args.yes
    project_id = gitlab_project_id(args.gitlab_project, gl_token)

    issues = fetch_github_issues(args.github_repo, gh_token, args.state)
    issues.sort(key=lambda i: i.number)

    skip_numbers: set[int] = set()
    for p in args.skip_map:
        try:
            data = json.load(open(p, "r", encoding="utf-8"))
            for entry in data.get("issues", []):
                n = entry.get("github_number")
                if isinstance(n, int):
                    skip_numbers.add(n)
        except FileNotFoundError:
            continue
        except Exception as ex:
            eprint(f"Warning: failed to read --skip-map {p}: {ex}")

    mapping: Dict[str, Any] = {
        "github_repo": args.github_repo,
        "gitlab_project": args.gitlab_project,
        "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "dry_run": dry_run,
        "issues": [],
    }

    eprint(f"GitHub issues to migrate: {len(issues)} (state={args.state})")
    if dry_run:
        eprint("Dry-run mode (no writes). Pass --yes to create issues/notes.")

    # Pre-create labels
    all_labels = sorted({l for i in issues for l in i.labels})
    for lbl in all_labels:
        ensure_gitlab_label(project_id, lbl, gl_token, dry_run)

    for issue in issues:
        if issue.number in skip_numbers:
            continue
        desc = format_issue_description(issue)
        eprint(f"#{issue.number} -> creating GitLab issue: {issue.title}")
        created = create_gitlab_issue(project_id, issue.title, desc, issue.labels, gl_token, dry_run)

        entry: Dict[str, Any] = {"github_number": issue.number, "github_url": issue.html_url}
        if created is not None:
            entry["gitlab_iid"] = created.get("iid")
            entry["gitlab_url"] = created.get("web_url")

        mapping["issues"].append(entry)

        if args.with_comments and created is not None:
            iid = int(created["iid"])
            comments = fetch_github_issue_comments(args.github_repo, issue.number, gh_token)
            for c in comments:
                note = format_comment_note(c)
                create_gitlab_note(project_id, iid, note, gl_token, dry_run)
                time.sleep(args.rate_limit_sleep)

        if args.close_when_closed and issue.state == "closed" and created is not None:
            close_gitlab_issue(project_id, int(created["iid"]), gl_token, dry_run)

        time.sleep(args.rate_limit_sleep)

    with open(args.out, "w", encoding="utf-8") as f:
        json.dump(mapping, f, indent=2, sort_keys=True)
        f.write("\n")

    eprint(f"Done. Wrote mapping to: {args.out}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
