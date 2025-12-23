# GitLab CI Reference

This repo includes a GitLab CI pipeline at `.gitlab-ci.yml` intended to replace the core GitHub Actions workflows:

- **CI** (`.github/workflows/ci.yml`): Go tests/lint, frontend lint/typecheck, Playwright smoke
- **Release** (`.github/workflows/release.yml`): Wails build artifacts for Linux/Windows and a tagged release

## Runner requirements

The pipeline is designed to run on Linux runners using Docker images.

Notes:

- The GitHub Actions pipeline currently uses a **macOS self-hosted runner** for its primary job. GitLab SaaS does not provide macOS runners by default—if you need macOS CI, register a macOS GitLab Runner and add a dedicated job/tag.
- Playwright runs in the official Playwright Docker image.

## Jobs

Core jobs:

- `go:lint` (golangci-lint)
- `go:test` (go test)
- `frontend:lint` (Biome)
- `frontend:typecheck` (TypeScript)
- `playwright:smoke` (Playwright E2E smoke)

Build jobs:

- `wails:build-linux` (produces `build/bin/*` artifacts)
- `wails:build-windows` (cross-compiles from Linux, produces `build/bin/*` artifacts)

Release job:

- `release` runs on tags that match `vX.Y.Z*` and creates a GitLab Release.

## Variables / secrets

If you migrate the AI agent workflows (Claude/Gemini/Codex) from GitHub Actions to GitLab, define equivalent GitLab CI/CD variables for the required API keys/tokens (e.g. `ANTHROPIC_API_KEY`).

## Migration tip: importing GitHub PRs/issues

If you want GitLab to preserve history (issues + PRs/MRs), prefer creating the GitLab project using GitLab’s GitHub import feature rather than pushing a new empty repo.
