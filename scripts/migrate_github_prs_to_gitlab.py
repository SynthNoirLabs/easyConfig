#!/usr/bin/env python3
"""
Migrate GitHub PRs into GitLab as:
- open PRs -> GitLab Merge Requests (pushes branches to GitLab)
- closed/merged PRs -> GitLab Issues (archival links)

Why: GitHub PR history cannot be migrated 1:1 to GitLab MRs without using platform importers.
This script preserves references and makes active work continue on GitLab.

Requires:
- GITHUB_TOKEN: GitHub token with repo access
- GITLAB_TOKEN: GitLab token with `api` scope

Usage:
  python3 scripts/migrate_github_prs_to_gitlab.py \
    --github-repo komod0/easyConfig \
    --gitlab-project komod0/easyconfig \
    --state all \
    --yes
"""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
import time
import urllib.parse
import urllib.request
from dataclasses import dataclass
from typing import Any, Dict, Iterable, List, Optional


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
class PullRequest:
    number: int
    title: str
    body: str
    state: str
    merged: bool
    html_url: str
    user_login: str
    created_at: str
    head_ref: str
    base_ref: str
    head_repo_full_name: str


def fetch_github_prs(repo: str, token: str, state: str) -> List[PullRequest]:
    url = f"{GITHUB_API}/repos/{repo}/pulls?state={urllib.parse.quote(state)}"
    prs: List[PullRequest] = []
    for item in paginate_github(url, token):
        prs.append(
            PullRequest(
                number=int(item["number"]),
                title=str(item["title"]),
                body=str(item.get("body") or ""),
                state=str(item["state"]),
                merged=bool(item.get("merged_at")),
                html_url=str(item["html_url"]),
                user_login=str(item["user"]["login"]),
                created_at=str(item["created_at"]),
                head_ref=str(item["head"]["ref"]),
                base_ref=str(item["base"]["ref"]),
                head_repo_full_name=str(item["head"]["repo"]["full_name"]) if item.get("head", {}).get("repo") else "",
            )
        )
    return prs


def gitlab_project_id(project_path: str, token: str) -> int:
    encoded = urllib.parse.quote(project_path, safe="")
    proj = http_json("GET", f"{GITLAB_API}/projects/{encoded}", gl_headers(token))
    return int(proj["id"])


def ensure_gitlab_label(project_id: int, label: str, token: str, dry_run: bool) -> None:
    if dry_run:
        return
    url = f"{GITLAB_API}/projects/{project_id}/labels"
    try:
        http_json("POST", url, gl_headers(token), {"name": label, "color": "#428BCA"})
    except HTTPError as ex:
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


def create_gitlab_mr(
    project_id: int,
    title: str,
    description: str,
    source_branch: str,
    target_branch: str,
    labels: List[str],
    token: str,
    dry_run: bool,
) -> Optional[Dict[str, Any]]:
    if dry_run:
        return None
    url = f"{GITLAB_API}/projects/{project_id}/merge_requests"
    payload: Dict[str, Any] = {
        "title": title,
        "description": description,
        "source_branch": source_branch,
        "target_branch": target_branch,
        "remove_source_branch": False,
        "squash": False,
    }
    if labels:
        payload["labels"] = ",".join(labels)
    return http_json("POST", url, gl_headers(token), payload)


def run_git(args: List[str]) -> Tuple[int, str]:
    proc = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
    return proc.returncode, proc.stdout


def push_branch_to_gitlab(source_ref: str, dest_branch: str) -> None:
    # Push from local origin/<branch> ref to gitlab branch
    code, out = run_git(["git", "push", "gitlab", f"{source_ref}:refs/heads/{dest_branch}"])
    if code != 0:
        raise RuntimeError(out.strip())


def format_mr_description(pr: PullRequest) -> str:
    header = (
        f"Migrated from GitHub PR: {pr.html_url}\n"
        f"Original author: @{pr.user_login}\n"
        f"Original created at: {pr.created_at}\n"
        f"\n---\n\n"
    )
    return header + (pr.body or "")


def format_pr_archive_issue(pr: PullRequest) -> Tuple[str, str]:
    title = f"[Archived PR] #{pr.number} {pr.title}"
    state = "merged" if pr.merged else pr.state
    body = (
        f"Migrated from GitHub PR: {pr.html_url}\n"
        f"Original author: @{pr.user_login}\n"
        f"Original created at: {pr.created_at}\n"
        f"Original state: {state}\n"
        f"Head: `{pr.head_repo_full_name}:{pr.head_ref}`\n"
        f"Base: `{pr.base_ref}`\n\n"
        f"---\n\n"
        f"{pr.body or ''}"
    )
    return title, body


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--github-repo", required=True)
    parser.add_argument("--gitlab-project", required=True)
    parser.add_argument("--state", default="open", choices=["open", "closed", "all"])
    parser.add_argument("--rate-limit-sleep", type=float, default=0.3)
    parser.add_argument("--yes", action="store_true")
    parser.add_argument("--out", default="migration-pr-map.json")
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

    prs = fetch_github_prs(args.github_repo, gh_token, args.state)
    prs.sort(key=lambda p: p.number)

    eprint(f"GitHub PRs to process: {len(prs)} (state={args.state})")
    if dry_run:
        eprint("Dry-run mode (no writes). Pass --yes to create issues/MRs.")

    # Labels used by this script
    ensure_gitlab_label(project_id, "github-pr", gl_token, dry_run)
    ensure_gitlab_label(project_id, "archived", gl_token, dry_run)

    mapping: Dict[str, Any] = {
        "github_repo": args.github_repo,
        "gitlab_project": args.gitlab_project,
        "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "dry_run": dry_run,
        "items": [],
    }

    for pr in prs:
        # Only create MRs for open PRs from same repo
        if pr.state == "open" and pr.head_repo_full_name.lower() == args.github_repo.lower():
            eprint(f"PR #{pr.number} -> MR from branch {pr.head_ref} -> {pr.base_ref}")
            # Ensure branch exists locally
            code, out = run_git(["git", "fetch", "origin", pr.head_ref])
            if code != 0:
                eprint(f"  fetch failed, archiving as issue: {out.strip()}")
            else:
                try:
                    push_branch_to_gitlab(f"origin/{pr.head_ref}", pr.head_ref)
                    mr = create_gitlab_mr(
                        project_id,
                        pr.title,
                        format_mr_description(pr),
                        pr.head_ref,
                        pr.base_ref,
                        ["github-pr"],
                        gl_token,
                        dry_run,
                    )
                    item: Dict[str, Any] = {"github_pr": pr.number, "github_url": pr.html_url, "type": "mr"}
                    if mr is not None:
                        item["gitlab_url"] = mr.get("web_url")
                        item["gitlab_iid"] = mr.get("iid")
                    mapping["items"].append(item)
                    time.sleep(args.rate_limit_sleep)
                    continue
                except Exception as ex:
                    eprint(f"  MR create failed, archiving as issue: {ex}")

        # Archive closed/merged PRs (and open fork PRs) as issues
        title, body = format_pr_archive_issue(pr)
        issue = create_gitlab_issue(
            project_id,
            title,
            body,
            ["github-pr", "archived"],
            gl_token,
            dry_run,
        )
        item = {"github_pr": pr.number, "github_url": pr.html_url, "type": "issue"}
        if issue is not None:
            item["gitlab_url"] = issue.get("web_url")
            item["gitlab_iid"] = issue.get("iid")
        mapping["items"].append(item)
        time.sleep(args.rate_limit_sleep)

    with open(args.out, "w", encoding="utf-8") as f:
        json.dump(mapping, f, indent=2, sort_keys=True)
        f.write("\n")

    eprint(f"Done. Wrote mapping to: {args.out}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

