## Why

This repo has six years of useful, hand-shaped CLI behavior but very little automated characterization around it. Before refactoring, we need comprehensive unit tests that pin down AWS SSO config parsing, cache token selection, credential output formatting, and command behavior without depending on real AWS accounts or a real home directory.

## What Changes

- Add a comprehensive Go unit test suite that preserves existing behavior across the core packages and CLI commands.
- Cover AWS profile parsing for default profiles, named profiles, legacy inline SSO keys, sso-session sections, profile override precedence, and current error messages for missing settings.
- Cover SSO cache token selection, including expired tokens, malformed cache files, unparsable expiry values, and the "run aws sso login" failure path.
- Cover command output contracts for `get --json`, human-readable `get`, `export`, `export-ps`, `helper`, `exec`, `set`, `list accounts`, and `list roles` where those commands can be exercised without real AWS calls.
- Add test seams or small internal helpers only where needed to replace AWS SDK calls, process execution, filesystem paths, or stdout/stderr capture in tests.
- Preserve user-visible command output, flag behavior, AWS config parsing semantics, SSO cache handling, credential_process JSON shape, and shell scripting compatibility.
- Non-goals: change AWS authentication semantics, alter credential storage behavior, introduce new CLI commands, modify credential output formats, or migrate away from Cobra/Viper/AWS SDK v2 as part of this change.

## Capabilities

### New Capabilities

- `behavior-preserving-unit-tests`: Comprehensive unit tests characterize existing aws-sso-creds behavior so future refactors can proceed with confidence.

### Modified Capabilities

None.

## Impact

- Affected code: `pkg/config`, `pkg/credentials`, `cmd/aws-sso-creds/**`, and any small testability hooks introduced near those packages.
- Affected APIs: no public CLI or Go API behavior should change.
- Affected dependencies: prefer the Go standard library for tests; add a dependency only if it materially improves maintainability.
- Affected systems: tests must not call real AWS APIs, read or write the user's actual `~/.aws` files, or require an existing SSO login.
