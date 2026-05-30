## Purpose

Document dependency hygiene expectations for keeping the module free of deprecated or unused AWS SDK dependencies while preserving CLI behavior.

## Requirements

### Requirement: Deprecated AWS SDK for Go v1 is absent
The module SHALL NOT depend on or import the deprecated `github.com/aws/aws-sdk-go` AWS SDK for Go v1 module.

#### Scenario: Module metadata excludes AWS SDK for Go v1
- **WHEN** the dependency graph is inspected through `go.mod` and `go.sum`
- **THEN** there MUST be no references to `github.com/aws/aws-sdk-go`

#### Scenario: Source code excludes AWS SDK for Go v1 imports
- **WHEN** repository Go source files are inspected
- **THEN** there MUST be no imports whose path starts with `github.com/aws/aws-sdk-go`

### Requirement: AWS SDK for Go v2 behavior is preserved
The system SHALL continue to use AWS SDK for Go v2 for AWS SSO API interactions without changing user-visible credential behavior.

#### Scenario: Credential commands retain output contracts
- **WHEN** unit tests exercise `get`, `export`, `export-ps`, `helper`, `exec`, `set`, `list accounts`, and `list roles`
- **THEN** those tests MUST continue to pass without requiring real AWS calls

#### Scenario: AWS SSO API calls use SDK v2 packages
- **WHEN** credential retrieval and list commands call AWS SSO APIs
- **THEN** those calls MUST use `github.com/aws/aws-sdk-go-v2` packages

### Requirement: Pulumictl dependency is absent
The module SHALL NOT depend on or import `github.com/pulumi/pulumictl` for version reporting.

#### Scenario: Module metadata excludes pulumictl
- **WHEN** the dependency graph is inspected through `go.mod` and `go.sum`
- **THEN** there MUST be no references to `github.com/pulumi/pulumictl`

#### Scenario: Source code excludes pulumictl imports
- **WHEN** repository Go source files are inspected
- **THEN** there MUST be no imports whose path starts with `github.com/pulumi/pulumictl`

### Requirement: Vers dependency is present for version reporting
The module SHALL use `github.com/jaxxstorm/vers` for fallback version calculation.

#### Scenario: Version command source uses vers
- **WHEN** the `version` command source is inspected
- **THEN** it MUST import or call `github.com/jaxxstorm/vers`
- **AND** it MUST NOT import `github.com/pulumi/pulumictl/pkg/gitversion`
