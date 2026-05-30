## 1. Test Infrastructure

- [x] 1.1 Add shared test helpers for temporary AWS home directories, config files, credentials files, SSO cache files, fake timestamps, and stdout capture.
- [x] 1.2 Add small package-local test seams for credential retrieval so tests can fake AWS SSO responses without network calls.
- [x] 1.3 Add small package-local test seams for command execution and path lookup so `exec` tests do not replace the test process.
- [x] 1.4 Add cleanup patterns for Viper globals, command flags, and replaceable seam functions so tests remain independent.

## 2. Config And Cache Tests

- [x] 2.1 Add `pkg/config` tests for default profile parsing with inline SSO settings.
- [x] 2.2 Add `pkg/config` tests for named profile parsing with inline SSO settings.
- [x] 2.3 Add `pkg/config` tests for `sso-session` merging and profile override precedence.
- [x] 2.4 Add `pkg/config` tests for missing profile and missing required setting errors.
- [x] 2.5 Add `pkg/config` tests for valid SSO cache token selection.
- [x] 2.6 Add `pkg/config` tests for expired tokens, non-matching start URLs, malformed JSON, unparseable expiration values, and the no-valid-cache-files error path.

## 3. Credential Retrieval Tests

- [x] 3.1 Add `pkg/credentials` tests proving successful retrieval uses the parsed account ID, role name, access token, and SSO region.
- [x] 3.2 Add `pkg/credentials` tests for missing cache directory errors.
- [x] 3.3 Add `pkg/credentials` tests proving AWS SSO API errors are wrapped with the existing retrieval context.

## 4. Credential Output Command Tests

- [x] 4.1 Add command tests for `get --json` that verify JSON field names and expiration representation.
- [x] 4.2 Add command tests for human-readable `get` output labels and credential fields.
- [x] 4.3 Add command tests for `export` POSIX shell assignment output.
- [x] 4.4 Add command tests for `export-ps` PowerShell assignment output.
- [x] 4.5 Add command tests for `helper` credential_process JSON keys, version, and RFC3339 expiration.

## 5. Side-Effect Command Tests

- [x] 5.1 Add command tests for `exec -- <command>` that capture the binary, args, and AWS environment variables.
- [x] 5.2 Add command tests for `set PROFILE` using temporary AWS config and credentials files.
- [x] 5.3 Add command tests for `list accounts` using fake SSO account data.
- [x] 5.4 Add command tests for `list roles ACCOUNT_ID` using fake SSO role data.

## 6. Verification

- [x] 6.1 Run `go test ./...` and fix failures without changing documented user-visible behavior.
- [x] 6.2 Run manual verification commands for output shape, including `go run ./cmd/aws-sso-creds --help`, `go run ./cmd/aws-sso-creds get --help`, `go run ./cmd/aws-sso-creds export --help`, and `go run ./cmd/aws-sso-creds helper --help`.
- [x] 6.3 Review tests for accidental real AWS access, real home directory reads/writes, real credential material, or brittle assertions unrelated to compatibility.
