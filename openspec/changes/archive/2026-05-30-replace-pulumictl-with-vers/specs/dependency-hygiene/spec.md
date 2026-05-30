## ADDED Requirements

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
