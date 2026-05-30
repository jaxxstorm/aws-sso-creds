## Context

`aws-sso-creds version` currently reads `pkg/version.Version`, which release builds populate through GoReleaser linker flags. When that value is empty, the command opens the current Git repository with `go-git` and asks `github.com/pulumi/pulumictl/pkg/gitversion` to calculate a semantic version from `HEAD`.

The Pulumi import is only present for version calculation, and the command's observable contract is small: print one version string to stdout or return an error if fallback calculation cannot be performed. The replacement library, `github.com/jaxxstorm/vers`, should take over that fallback calculation without changing credential-related behavior.

## Goals / Non-Goals

**Goals:**

- Use `github.com/jaxxstorm/vers` for fallback version calculation.
- Remove direct `github.com/pulumi/pulumictl` imports and module requirements.
- Preserve `aws-sso-creds version` output shape and linked-version precedence.
- Keep the change isolated to version reporting and dependency metadata.

**Non-Goals:**

- Changing GoReleaser linker flag names or release version injection.
- Changing AWS profile parsing, SSO cache lookup, credential retrieval, or credential output.
- Introducing new version command flags or changing command names.

## Decisions

### Decision: Keep linked version precedence

The command will continue to print `pkg/version.Version` when it is non-empty. This preserves release build behavior and avoids invoking Git inspection in normal packaged binaries.

Alternative considered: always ask `vers` for the version. That would make release binaries depend on local Git metadata at runtime and could change output for installed artifacts.

### Decision: Use `vers` only for fallback calculation

`vers` should replace the Pulumi-specific fallback version calculation. The command should no longer import `github.com/pulumi/pulumictl/pkg/gitversion`; any direct `go-git` usage should also be removed if `vers` provides the repository inspection entry point itself.

Alternative considered: wrap the existing Pulumi API behind a local helper first. That adds indirection without achieving the dependency cleanup the change is for.

### Decision: Test behavior, not Git internals

Tests should cover linked-version output and fallback output/error behavior through a narrow seam around version calculation. The tests should not require real AWS calls, home directories, or credential fixtures because this command does not touch credential state.

Alternative considered: rely on end-to-end command execution inside the repository. That is useful as a smoke test, but unit coverage should keep fallback behavior deterministic.

## Risks / Trade-offs

- `vers` may format fallback versions slightly differently than Pulumi's helper. Mitigation: add tests for the expected output contract and verify the command manually before archiving.
- Removing `go-git` from direct dependencies may reshuffle indirect dependencies. Mitigation: run `go mod tidy` and inspect `go.mod`/`go.sum` for the intended Pulumi removal.
- The old command long description says `pulumictl`. Mitigation: update the description while preserving command name, short help, and output behavior.
