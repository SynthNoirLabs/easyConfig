# easyConfig Repository Review

## 1. Executive Summary

`easyConfig` is an ambitious Wails desktop application that tries to solve a real operational problem: every AI coding tool ships its own configuration format, storage location, and extension mechanism, and the project provides a single place to discover, edit, compare, back up, and extend those configs. The core architectural choice is sound. The Go backend uses a provider pattern in `pkg/config`, the React frontend uses a small Context-based state layer over Wails RPC bindings, and cross-cutting concerns such as MCP injection, schema fetching, profiles, workflows, and file watching are separated into dedicated packages. For a relatively small codebase, it already covers a broad surface area well.

The codebase is not production-ready yet in its current form, mainly because the trust boundaries are still too soft for a configuration manager that touches secrets and writes arbitrary files on the user machine. Several critical concerns are backend path validation, secret handling via command-line arguments, incomplete platform/path normalization around Claude and MCP integration, and CI/build ergonomics that still require bootstrap workarounds such as creating `frontend/dist` before `go test ./...`. None of these problems require a rewrite, but they do require deliberate hardening before the application should be positioned as a secure daily driver.

## 2. Architectural Strengths

### Clear separation between backend orchestration and frontend presentation

- `main.go` keeps the Wails bootstrap small and binds a single `App` service object.
- `app.go` acts as a thin application boundary instead of mixing UI concerns into package code.
- `frontend/src/context/ConfigContext.tsx` centralizes Wails calls and file-change subscriptions, so most components remain presentation-focused.

This is a good fit for Wails. The bindings are simple, easy to regenerate, and easy to test in React by mocking `window.go` and `window.runtime`.

### Provider pattern is the right abstraction for the domain

The `Provider` interface in `pkg/config/types.go` is a strong extension point:

```go
type Provider interface {
    Name() string
    Discover(projectPath string) ([]Item, error)
    Create(scope Scope, projectPath string) (string, error)
    CheckStatus() ProviderStatus
    BinaryName() string
    VersionArgs() []string
}
```

Strengths:

- adding a new agent/tool is mostly localized to a new `provider_*.go` file,
- discovery logic is decoupled from UI and persistence,
- health/status checks are first-class instead of bolted on later,
- the service can fan out provider discovery concurrently.

For a tool that will keep growing as the AI CLI ecosystem changes, this is the right direction.

### Sensible package decomposition

The repository does a good job isolating responsibilities:

- `pkg/config`: discovery, profiles, import/export, provider definitions,
- `pkg/mcp`: MCP config injection,
- `pkg/watcher`: file change monitoring,
- `pkg/schema`: schema fetching,
- `pkg/workflows`: workflow generation and secret helpers,
- `pkg/util/paths`: path normalization and home/config directory logic.

That decomposition makes future hardening feasible because the security-sensitive seams are easy to find.

### Restrictive file permissions are already part of the design

Multiple write paths use `0600`, including config saves, MCP writes, and exported data. That is exactly the right default for agent configuration files, which often contain API keys or tokens. This is one of the better decisions in the codebase.

### Test coverage is broad for the repository size

The project already has:

- Go unit tests across `pkg/config`, `pkg/mcp`, `pkg/schema`, `pkg/workflows`, `pkg/watcher`, and `pkg/install`,
- frontend tests for context and several components,
- integration tests under `tests/integration`.

That is a strong foundation. The main gap is not absence of tests; it is that some of the highest-risk paths are not yet under security-focused tests.

## 3. Critical Vulnerabilities & Bugs

Ranked from highest severity to lowest.

### Critical: backend file operations trust UI-supplied paths too much

`pkg/config/service.go` exposes:

- `ReadConfig(path string)`
- `SaveConfig(path, content string)`
- `DeleteConfig(path string)`

These ultimately read, write, or delete arbitrary filesystem paths without verifying that the path belongs to a discovered config item, a known profile root, or an allowed scope. In a Wails desktop app the frontend is local, but it is still an untrusted boundary relative to the backend. Any XSS-style issue, malicious extension, or future plugin system could turn this into arbitrary local file read/write/delete.

Impact:

- overwrite arbitrary files writable by the current user,
- delete unrelated local files,
- read secrets outside the config domain,
- apply profiles to unsafe paths persisted earlier.

Recommendation:

1. Introduce a backend path authorization layer.
2. Resolve paths with `filepath.Clean` and `filepath.EvalSymlinks` where possible.
3. Allow only:
   - files previously discovered by registered providers,
   - files created within controlled application directories,
   - explicit import/export destinations selected by a native file picker flow.

Example direction:

```go
func (s *DiscoveryService) validateManagedPath(path string) error {
    cleaned := filepath.Clean(path)
    if !filepath.IsAbs(cleaned) {
        return fmt.Errorf("path must be absolute")
    }
    if !s.isKnownManagedPath(cleaned) {
        return fmt.Errorf("path is not a managed configuration file")
    }
    return nil
}
```

### High: repository secrets are passed to `gh` using `--body`

`pkg/workflows/secrets.go` currently runs:

```go
exec.Command("gh", "secret", "set", name, "--body", value)
```

Passing the secret value on the command line risks exposure through:

- local process inspection,
- shell history in future refactors,
- debug logs or command tracing,
- CI telemetry.

Recommendation: send the secret through stdin instead.

Safer pattern:

```go
cmd := exec.Command("gh", "secret", "set", name, "--body-file", "-")
cmd.Stdin = strings.NewReader(value)
```

### High: MCP injection uses inconsistent Claude Desktop path assumptions

`app.go` contains a long comment trail acknowledging uncertainty, then still writes to:

```go
filepath.Join(homeDir, ".claude", "claude_desktop_config.json")
```

That is a production risk, not just a style issue:

- wrong path on macOS/Windows means silent failure or split-brain configs,
- discovery and mutation can drift apart,
- users may believe an MCP package was installed when the real client never reads the file.

Recommendation:

- move Claude Desktop path resolution into `pkg/util/paths` or the Claude provider,
- use one canonical resolver shared by discovery and injection,
- unit test Linux/macOS/Windows path derivation.

### High: import and profile application flows can replay unsafe paths

`pkg/config/import.go` stores imported profile items using a placeholder path, which is good for remote imports, but local profiles in `profiles.go` preserve original absolute paths and `ApplyProfile` writes them back without re-validating that the destination is still safe or expected.

Risks:

- stale paths pointing to unrelated files after machine changes,
- symlink abuse between backup and restore operations,
- restoring profiles across environments with different directory ownership or semantics.

Recommendation:

- validate each profile item path before backup or write,
- bind profile items to provider + scope + relative identifier instead of raw absolute path wherever possible,
- require confirmation when applying a profile containing paths outside the current OS-specific provider roots.

### Medium: CI/bootstrap is brittle because Go embed requires a frontend artifact

`main.go` embeds `all:frontend/dist`, which makes the root package fail under `go test ./...` in a fresh clone unless CI or the developer pre-creates the directory. The CI workaround is:

```yaml
- name: Ensure dist exists for Go embed
  run: mkdir -p frontend/dist && touch frontend/dist/.keep
```

This is operationally fragile:

- local contributors hit false-negative test failures,
- root-package tests depend on generated frontend state,
- tooling feels broken before bootstrap completes.

Recommendation:

- commit a tiny placeholder file under `frontend/dist/.keep`, or
- move embed behind build tags / a runtime asset server strategy for tests.

### Medium: import from URL lacks provenance and size controls

`ImportProfilesFromURL` uses `http.Get` directly and accepts any successful response body into memory before unmarshalling. Missing controls include:

- timeout-bound custom client,
- size limit via `io.LimitReader`,
- scheme allowlist,
- optional checksum/signature validation for shared profile bundles.

This is not immediately RCE, but it is a weak ingestion point for a tool that may later be marketed as a secure config manager.

### Medium: agentic CI workflows have overbroad trust and side effects

The agent workflows are innovative, but several of them combine write permissions, external model execution, and issue/PR-triggered automation. Examples include:

- `agent-interactive.yml` with `contents: write`, `issues: write`, `pull-requests: write`,
- `gemini-review.yml` minting tokens and allowing PR review actions,
- `gemini-agent.yml` performing automated review/triage/chat directly from GitHub events.

Risks:

- prompt injection via issue comments or PR descriptions,
- accidental disclosure of repo context to third-party model providers,
- automated approval/review comments based on manipulated diffs,
- operational noise or runaway automation loops.

Recommendation:

- isolate agent workflows to least privilege,
- separate read-only review from write-capable fix workflows,
- gate write operations behind labels from trusted maintainers only,
- disable agent execution on forks unless a maintainer explicitly re-triggers.

## 4. Code Smells & Tech Debt

### Hardcoded knowledge is creeping into `app.go`

`app.go` has grown into a mixed orchestration layer containing:

- Wails-bound methods,
- marketplace aggregation,
- MCP installation logic,
- profile import/export endpoints,
- path commentary that belongs in dedicated utilities.

This is not broken yet, but it is heading toward a God object. The next inflection point should split it into narrow services and keep `App` as a facade.

### Path utilities are not consistently reused

The repository explicitly documents “use `pkg/util/paths`”, but not all code follows that rule. The Claude injection logic is the clearest example. This creates subtle cross-platform drift, which is especially dangerous in a desktop app that manages user files.

### Frontend state model is clean but still centralized and coarse

`ConfigContext.tsx` is intentionally simple, but all config list loading, read, save, delete, and change-notification behavior lives in one context. That is acceptable now, but if profiles, marketplace, docs, workflows, and health continue growing, this single context will become an accidental global store.

Recommendation:

- keep `ConfigContext` for config CRUD,
- add dedicated contexts or query hooks for health, docs, marketplace, and workflows.

### Type safety is decent, but model drift is still a risk

The frontend correctly imports Wails-generated models, which is good, but domain logic also introduces local composite types like:

```ts
type SelectableItem = config.Item & { initialLine?: number };
```

This is fine in isolation, but the broader risk is that generated Wails models and hand-authored UI assumptions can diverge quietly. The project should prefer thin view-model adapters close to component boundaries instead of mutating generated types in place across the UI.

### Error handling favors “continue quietly” in too many places

Examples:

- provider discovery logs errors and continues,
- profile listing skips unreadable/bad files,
- preview/diff generation skips files it cannot read.

That behavior is user-friendly, but the UI needs structured warnings so users understand partial failure. Silent degradation is acceptable for discovery; it is less acceptable for profile application and backup flows.

### Accessibility and keyboard support are present but incomplete

The frontend shows good intent with keyboard shortcuts and modal structure, but there is not enough evidence of systematic accessibility work:

- no visible a11y audit tooling in CI,
- uncertain ARIA coverage for custom components,
- icon-only settings button depends on `title`, which is weaker than an explicit accessible label.

This is not catastrophic for an internal utility, but it is a production-readiness gap.

### Tests do not yet focus on abuse cases

The repository has many tests, but high-value missing cases include:

- rejecting unsafe config paths,
- profile apply with symlinked destinations,
- secrets manager ensuring values are never passed as process args,
- import from URL enforcing size/time limits,
- Windows/macOS/Linux path resolution for Claude Desktop and other providers.

## 5. Actionable Recommendations

### A. Tighten backend trust boundaries first

If only one hardening track is prioritized, make it this one.

Recommended backend flow:

```text
React UI
   |
   v
Wails method in App
   |
   v
Input validation / path authorization
   |
   +--> provider-managed config roots
   +--> easyConfig-owned data roots
   +--> explicit import/export destinations
   |
   v
pkg/config or pkg/mcp write operation
```

Concretely:

1. add a `validateManagedPath` function in `pkg/config/service.go`,
2. call it from `ReadConfig`, `SaveConfig`, `DeleteConfig`, `ApplyProfile`, and backup restore flows,
3. build tests that prove arbitrary files outside managed roots are rejected.

### B. Refactor secret writes to stdin-based transport

Suggested replacement for `pkg/workflows/secrets.go`:

```go
func (s *SecretsManager) SetRepositorySecret(name, value string) error {
    if _, err := exec.LookPath("gh"); err != nil {
        return fmt.Errorf("GitHub CLI (gh) is not installed or not in PATH")
    }

    cmd := exec.Command("gh", "secret", "set", name, "--body-file", "-")
    cmd.Stdin = strings.NewReader(value)

    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("failed to set secret: %s: %s", err, string(output))
    }
    return nil
}
```

Add a unit test that asserts the command args do not contain the secret value.

### C. Centralize provider path resolution

Introduce a small resolver layer:

```go
type ProviderPaths interface {
    ConfigPaths(home string, project string) []string
}
```

Then:

- keep path definitions near each provider,
- export canonical helpers for mutation code such as MCP injection,
- reuse the same functions in discovery, creation, health checks, and editing.

This reduces the chance that “discover path A, mutate path B” bugs recur as more providers are added.

### D. Make CI reflect real developer bootstrap

The current CI already documents the fix implicitly: create `frontend/dist` before Go tests and run `npm ci` before frontend commands. Make that experience explicit for local contributors too.

Recommended improvements:

1. include a checked-in placeholder file in `frontend/dist/`,
2. add a root bootstrap command to README/Taskfile that clearly states frontend deps are required for frontend lint/test/build,
3. consider a lightweight preflight target:

```bash
task setup && mkdir -p frontend/dist && touch frontend/dist/.keep
```

### E. Harden agentic CI

A safe pattern is to split agent workflows into three classes:

```text
Class 1: Read-only analysis
- PR review comments
- issue triage suggestions
- no repo write token

Class 2: Draft remediation
- artifacts or patch suggestions
- gated by maintainer label
- no direct merge authority

Class 3: Write-capable automation
- explicit maintainer dispatch only
- short-lived app token
- isolated branch / sandbox
```

Operationally:

- never grant `contents: write` to comment-triggered workflows by default,
- strip untrusted prompt material before passing issue comments to external models,
- log exactly what repository data is sent to third-party model providers,
- require maintainer approval before posting reviews or code changes back.

### F. Improve observability for partial failures

Instead of silently continuing on malformed files, return structured warnings:

```go
type DiscoverResult struct {
    Items    []Item           `json:"items"`
    Warnings []ProviderWarning `json:"warnings"`
}
```

That would let the frontend surface:

- “7 providers scanned”
- “2 providers partially failed”
- “1 profile file skipped because JSON was invalid”

This keeps the forgiving UX without sacrificing debuggability.

### G. Add a small production-readiness checklist

Before calling the app production-ready, the project should be able to answer “yes” to all of the following:

- Are backend file writes constrained to managed paths?
- Are secrets never passed via process args or logs?
- Are provider paths derived from one canonical cross-platform source?
- Can `go test ./...` succeed in a fresh clone without hidden prerequisites?
- Are agentic CI workflows least-privilege and fork-safe?
- Are high-risk filesystem flows covered by abuse-case tests?

If those items are completed, the existing architecture is strong enough to scale.
