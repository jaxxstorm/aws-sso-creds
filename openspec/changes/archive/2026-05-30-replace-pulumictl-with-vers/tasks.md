## 1. Dependency Update

- [x] 1.1 Add `github.com/jaxxstorm/vers` as the version fallback dependency.
- [x] 1.2 Remove `github.com/pulumi/pulumictl` imports from the version command.
- [x] 1.3 Run `go mod tidy` so `go.mod` and `go.sum` remove unused Pulumi entries and retain only required versioning dependencies.

## 2. Version Command Implementation

- [x] 2.1 Update `cmd/aws-sso-creds/version/cli.go` to call `vers` when `pkg/version.Version` is empty.
- [x] 2.2 Preserve linked-version precedence so non-empty `pkg/version.Version` prints directly without fallback calculation.
- [x] 2.3 Update the command long description to refer to `aws-sso-creds` instead of `pulumictl`.
- [x] 2.4 Keep stdout behavior to exactly one version string followed by a newline on success.

## 3. Tests

- [x] 3.1 Add or update version command unit tests for linked-version output.
- [x] 3.2 Add or update version command unit tests for fallback version calculation using a deterministic seam or fixture.
- [x] 3.3 Add or update version command unit tests for fallback calculation errors.
- [x] 3.4 Add a test or static assertion that version help text mentions `aws-sso-creds` and not `pulumictl`.

## 4. Verification

- [x] 4.1 Run `go test ./...`.
- [x] 4.2 Run `go run cmd/aws-sso-creds/main.go version` and confirm it prints a single version line.
- [x] 4.3 Inspect `go.mod`, `go.sum`, and source imports to confirm `github.com/pulumi/pulumictl` is absent.
