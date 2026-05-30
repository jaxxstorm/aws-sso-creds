## 1. Dependency Audit

- [x] 1.1 Search the repository for `github.com/aws/aws-sdk-go` imports or references.
- [x] 1.2 Confirm current AWS-facing source code imports `github.com/aws/aws-sdk-go-v2` packages for SSO behavior.

## 2. Module Cleanup

- [x] 2.1 Remove the direct `github.com/aws/aws-sdk-go` dependency from `go.mod`.
- [x] 2.2 Run `go mod tidy` to remove unused v1 SDK checksums from `go.sum`.
- [x] 2.3 Review `go.mod` and `go.sum` diffs to ensure the cleanup stays scoped to unused dependency removal.

## 3. Verification

- [x] 3.1 Re-run repository search and confirm `go.mod`, `go.sum`, and source files contain no `github.com/aws/aws-sdk-go` references.
- [x] 3.2 Run `go test ./...` to verify existing behavior-preserving tests still pass.
- [x] 3.3 Run `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run` to verify lint remains clean.
- [x] 3.4 Manually verify a representative CLI command such as `go run ./cmd/aws-sso-creds get --help` still runs without changing output shape.
