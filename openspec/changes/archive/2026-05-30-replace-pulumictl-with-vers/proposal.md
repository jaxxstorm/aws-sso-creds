## Why

The `version` command currently imports Pulumi's `pulumictl` only to calculate a fallback version when linker flags have not set one. That versioning logic now lives in `github.com/jaxxstorm/vers`, so this repository can depend on the smaller purpose-built module instead of a tool-specific package from another project.

## What Changes

- Replace `github.com/pulumi/pulumictl/pkg/gitversion` usage in `cmd/aws-sso-creds/version` with `github.com/jaxxstorm/vers`.
- Remove `github.com/pulumi/pulumictl` from module dependencies and checksums.
- Preserve the existing `aws-sso-creds version` command behavior: print the linked version when present, otherwise calculate a version from the current Git repository and print it to stdout.
- Fix the `version` command long description so it refers to `aws-sso-creds` rather than `pulumictl`.
- No changes to AWS config parsing, SSO cache handling, credential_process output, credential storage, or shell scripting output.

## Capabilities

### New Capabilities

- `cli-version-reporting`: Defines the observable behavior of the `version` command and its fallback version calculation.

### Modified Capabilities

- `dependency-hygiene`: Adds a requirement that the module must not import or depend on `github.com/pulumi/pulumictl` once versioning has moved to `github.com/jaxxstorm/vers`.

## Impact

- Affected code: `cmd/aws-sso-creds/version/cli.go`, `pkg/version`, and version command tests.
- Affected dependencies: `go.mod` and `go.sum` should remove `github.com/pulumi/pulumictl` and include `github.com/jaxxstorm/vers`.
- Affected CLI commands: only `aws-sso-creds version`; expected output remains one version string followed by a newline.
- Non-goals: changing release-time linker flags, AWS authentication semantics, credential output formats, SSO cache lookup behavior, or any credential persistence behavior.
