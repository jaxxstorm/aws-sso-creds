## Context

The project already uses AWS SDK for Go v2 in source code for AWS SSO configuration loading and SSO API calls. The deprecated AWS SDK for Go v1 remains as a direct module dependency in `go.mod` and appears in `go.sum`, but current source imports do not reference it.

This change should be a module dependency cleanup, not a rewrite of credential retrieval.

## Goals / Non-Goals

**Goals:**

- Remove `github.com/aws/aws-sdk-go` from the module graph when it is unused.
- Keep AWS SSO credential retrieval on AWS SDK for Go v2.
- Verify no Go source imports the v1 SDK.
- Preserve existing CLI behavior, command output, tests, and lint behavior.

**Non-Goals:**

- Do not change AWS profile parsing.
- Do not change SSO cache lookup or validation.
- Do not change credential_process, shell export, PowerShell export, or human-readable output formats.
- Do not change AWS authentication semantics or credential storage.
- Do not introduce new AWS SDK abstractions unless removal of v1 exposes an actual compile-time need.

## Decisions

### Decision: Treat this as dependency removal, not SDK migration

The codebase has already migrated runtime AWS calls to AWS SDK for Go v2. Implementation should first prove whether v1 is only stale module metadata, then remove it with `go mod tidy`.

Alternative considered: proactively refactor all AWS-facing code. That would increase risk without addressing the actual deprecated dependency reference.

### Decision: Verify imports before and after module cleanup

Use repository-wide search for `github.com/aws/aws-sdk-go` and `github.com/aws/aws-sdk-go/` to ensure no source file imports v1. After tidying, repeat the search against source and module files so the deprecated module is gone.

Alternative considered: rely only on `go mod tidy`. Tidy can remove unused modules, but explicit search gives a clearer guardrail for future review.

### Decision: Preserve existing AWS SDK v2 module versions unless Go chooses otherwise

Do not intentionally upgrade AWS SDK for Go v2 modules as part of this change. If `go mod tidy` adjusts indirect dependencies mechanically, review the diff and keep it scoped to what removal requires.

Alternative considered: upgrade all AWS SDK v2 modules at the same time. That is a separate dependency update with wider test surface.

## Risks / Trade-offs

- `go mod tidy` may remove or adjust more checksums than just v1 -> Review `go.mod` and `go.sum` diffs and keep only module-graph changes tied to the cleanup.
- A hidden v1 import could exist in generated or build-tagged code -> Use `rg` across the repository and rely on `go test ./...` to compile normal packages.
- CI may exercise lint in addition to tests -> Run the same golangci-lint v2 command used locally for the config fix.
