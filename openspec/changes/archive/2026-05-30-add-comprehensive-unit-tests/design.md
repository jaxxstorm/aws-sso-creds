## Context

`aws-sso-creds` is a small Go CLI with behavior that has stayed useful for years: read AWS SSO profile config, find a cached access token, call AWS SSO for role credentials, and print or apply those credentials in several command-specific formats. The code currently mixes command construction, stdout writes, filesystem reads, AWS SDK calls, and process replacement in ways that make end-to-end behavior hard to unit test safely.

The goal is not to redesign the tool before refactoring. The goal is to add characterization tests that document current behavior and create just enough seams to run those tests without touching real AWS state.

## Goals / Non-Goals

**Goals:**

- Cover the behavior users depend on before larger refactors happen.
- Exercise `pkg/config` with temporary AWS config and SSO cache fixtures.
- Exercise credential retrieval with fake AWS SSO clients instead of network calls.
- Exercise CLI output contracts for shell, PowerShell, credential_process, JSON, human-readable, list, set, and exec workflows.
- Keep testability changes small, local, and compatible with the existing Cobra/Viper package layout.

**Non-Goals:**

- Do not change command names, flags, output formats, error wording, or AWS auth semantics.
- Do not migrate command packages to a new CLI framework.
- Do not replace AWS SDK v2 or introduce integration tests against live AWS accounts.
- Do not require real `~/.aws` files, real SSO login state, or real credentials.

## Decisions

### Decision: Add package-local test seams instead of a broad dependency injection framework

Use small function variables or narrow interfaces at the boundaries that are currently hard to test: credential retrieval, SSO API calls, stdout/stderr writers, process execution, and command path lookup. Keep these seams close to the packages that need them and reset them in tests.

Alternative considered: rewrite command constructors around a shared application struct. That may be useful later, but it is too much architectural movement for a behavior-preserving test change.

### Decision: Keep filesystem behavior real but redirected to temporary homes

Tests for AWS config parsing, cache parsing, and `set` should write realistic fixture files under `t.TempDir()` and pass `--home-directory` or direct homedir arguments. This preserves the file format behavior that matters while avoiding user state.

Alternative considered: mock all file reads and writes. That would make tests less representative of the actual AWS file layouts this tool exists to handle.

### Decision: Fake AWS SDK calls at the narrowest useful boundary

Credential retrieval and list commands should be able to use fake SSO clients that return deterministic role credentials, accounts, roles, regions, and errors. Tests should assert request inputs such as account ID, role name, access token, and SSO region where practical.

Alternative considered: use AWS SDK stub middleware or local HTTP servers. That adds complexity without improving confidence for this CLI's behavior.

### Decision: Treat command output as a compatibility contract

Tests should assert stable labels, environment variable names, JSON keys, and table headers. Where timestamps or colored output make byte-for-byte matching brittle, tests can parse JSON or assert important substrings, but they should still protect the visible contract users script against.

Alternative considered: only test internal helper functions. That would miss the user-facing behavior this change is meant to preserve.

### Decision: Avoid real process replacement in `exec` tests

The `exec` command currently ends by calling `syscall.Exec`. Introduce a small replaceable execution function so tests can capture the binary, args, and environment instead of replacing the test process.

Alternative considered: run a child process fixture command. That is useful for integration coverage, but unit tests should stay fast and deterministic.

## Risks / Trade-offs

- Test seams accidentally alter production behavior -> Keep default seam values wired to the current functions and add tests around the public command behavior after seams are introduced.
- Output tests become brittle around harmless formatting changes -> Assert the parts users rely on, such as field names, environment variable names, table headers, and credential_process keys.
- Fake AWS clients drift from AWS SDK shapes -> Keep fake interfaces close to the methods currently used and use AWS SDK response types where that does not complicate tests.
- Viper global state leaks between command tests -> Reset relevant Viper keys and command flags in test cleanup, and prefer constructing fresh commands per test.
- `set` tests may expose existing defects in file creation or error handling -> Characterize current behavior first, then decide separately whether any discovered bug deserves its own behavior-changing OpenSpec change.

## Migration Plan

Add tests and test seams incrementally by package. Run `go test ./...` after each cluster of changes. Because this is a test-only behavior-preserving change, rollback is simply reverting the test files and any small testability hooks.

## Open Questions

- Should any currently surprising behavior discovered while writing tests be preserved exactly, or should it become a separate follow-up fix? Default answer: preserve it here and document follow-up fixes separately.
