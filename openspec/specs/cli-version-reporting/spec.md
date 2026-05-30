## Purpose

Document the observable behavior of the `aws-sso-creds version` command and its fallback Git-based version calculation.

## Requirements

### Requirement: Version command reports linked version
The system SHALL make the `version` command print the linked `pkg/version.Version` value when that value is non-empty.

#### Scenario: Linked version is available
- **WHEN** `aws-sso-creds version` runs with a non-empty linked version value
- **THEN** stdout MUST contain that version followed by a newline
- **AND** fallback Git version calculation MUST NOT be required

### Requirement: Version command calculates fallback version
The system SHALL calculate and print a fallback version from the current Git repository when no linked version value is present.

#### Scenario: Linked version is absent
- **WHEN** `aws-sso-creds version` runs with an empty linked version value inside a readable Git repository
- **THEN** stdout MUST contain the calculated version followed by a newline

#### Scenario: Fallback version cannot be calculated
- **WHEN** `aws-sso-creds version` runs with an empty linked version value and fallback version calculation fails
- **THEN** the command MUST return an error
- **AND** the command MUST NOT print a credential, token, or AWS profile value

### Requirement: Version command help names aws-sso-creds
The system SHALL describe the `version` command as reporting the version of `aws-sso-creds`.

#### Scenario: Version help is displayed
- **WHEN** help text for `aws-sso-creds version` is rendered
- **THEN** the long description MUST refer to `aws-sso-creds`
- **AND** the long description MUST NOT refer to `pulumictl`
