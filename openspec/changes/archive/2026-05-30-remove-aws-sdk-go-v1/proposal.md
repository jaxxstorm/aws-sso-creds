## Why

`github.com/aws/aws-sdk-go` is the deprecated AWS SDK for Go v1, while this project already uses AWS SDK for Go v2 for its SSO calls. Removing the v1 dependency keeps the dependency graph current and avoids carrying an obsolete AWS SDK that is not part of the runtime behavior.

## What Changes

- Remove the direct `github.com/aws/aws-sdk-go` dependency from `go.mod`.
- Run module tidying so `go.sum` no longer contains AWS SDK for Go v1 checksums when they are unused.
- Keep all AWS SSO API calls on AWS SDK for Go v2 packages.
- Preserve existing CLI behavior and output for `get`, `export`, `export-ps`, `exec`, `helper`, `set`, `list accounts`, and `list roles`.
- Non-goals: change AWS profile parsing, SSO cache handling, AWS authentication semantics, credential storage, credential_process output, or shell scripting compatibility.

## Capabilities

### New Capabilities
- `dependency-hygiene`: Ensure the module avoids deprecated or unused AWS SDK dependencies while preserving runtime behavior.

### Modified Capabilities
None. This is a dependency cleanup with no intended spec-level behavior change.

## Impact

- Affected files: `go.mod`, `go.sum`, and any code or tests that still import AWS SDK for Go v1 if found during implementation.
- Affected dependencies: remove `github.com/aws/aws-sdk-go`; retain AWS SDK for Go v2 modules already used by credential retrieval and list commands.
- Affected CLI commands: all commands should remain behaviorally unchanged.
- User-visible output: no intended changes.
