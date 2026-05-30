## ADDED Requirements

### Requirement: Config parsing behavior is characterized
The test suite SHALL preserve AWS SSO configuration parsing behavior for supported profile layouts and failure modes.

#### Scenario: Default profile is parsed
- **WHEN** a fixture home directory contains a default AWS config section with inline SSO settings
- **THEN** tests MUST verify that the default profile resolves the expected start URL, SSO region, account ID, and role name

#### Scenario: Named profile is parsed
- **WHEN** a fixture home directory contains a named profile section with inline SSO settings
- **THEN** tests MUST verify that the named profile resolves the expected start URL, SSO region, account ID, and role name

#### Scenario: Shared sso-session settings are merged
- **WHEN** a profile references an `sso-session` section
- **THEN** tests MUST verify that SSO settings are read from the session and profile-specific account and role settings are preserved

#### Scenario: Profile settings override shared sso-session settings
- **WHEN** both a profile and its referenced `sso-session` section define the same SSO setting
- **THEN** tests MUST verify that the profile value wins

#### Scenario: Missing profile returns existing error
- **WHEN** the requested profile section does not exist
- **THEN** tests MUST verify that the existing "unable to find profile" error behavior is preserved

#### Scenario: Missing required SSO settings return existing errors
- **WHEN** a profile or session omits a required SSO setting
- **THEN** tests MUST verify that the existing missing-setting error behavior is preserved

### Requirement: SSO cache token selection is characterized
The test suite SHALL preserve SSO cache token selection behavior without reading the user's real AWS SSO cache.

#### Scenario: Valid matching cache token is selected
- **WHEN** fixture cache files include an unexpired token whose start URL matches the selected profile
- **THEN** tests MUST verify that the matching access token is returned

#### Scenario: Expired tokens are ignored
- **WHEN** fixture cache files include a matching token whose expiration is in the past
- **THEN** tests MUST verify that the token is ignored and the no-valid-cache-files behavior is preserved if no later valid token exists

#### Scenario: Non-matching start URLs are ignored
- **WHEN** fixture cache files include tokens for other SSO start URLs
- **THEN** tests MUST verify that those tokens are ignored

#### Scenario: Malformed JSON returns existing error
- **WHEN** a cache file contains malformed JSON
- **THEN** tests MUST verify that the existing JSON parsing error behavior is preserved

#### Scenario: Unparseable expiration is skipped
- **WHEN** a matching cache file has an expiration value that cannot be parsed
- **THEN** tests MUST verify that the cache file is skipped and later valid files may still be used

### Requirement: Credential retrieval is testable without AWS
The test suite SHALL preserve credential retrieval behavior while replacing real AWS SSO API calls with deterministic fakes.

#### Scenario: Successful retrieval uses parsed config and token
- **WHEN** fixture config and cache data resolve to a valid SSO role request
- **THEN** tests MUST verify that credential retrieval asks for the expected account ID, role name, access token, and SSO region

#### Scenario: Missing cache directory preserves current error path
- **WHEN** the fixture home directory does not contain an SSO cache directory
- **THEN** tests MUST verify that credential retrieval returns the existing login-oriented cache error

#### Scenario: AWS API errors are wrapped
- **WHEN** the fake SSO client returns an error while retrieving role credentials
- **THEN** tests MUST verify that credential retrieval returns the existing AWS retrieval error context

### Requirement: CLI output contracts are characterized
The test suite SHALL preserve user-visible output for credential-producing commands.

#### Scenario: Get command JSON output is stable
- **WHEN** `get --json` receives deterministic fake credentials
- **THEN** tests MUST verify the JSON field names and expiration representation stay compatible

#### Scenario: Get command human output is stable
- **WHEN** `get` receives deterministic fake credentials
- **THEN** tests MUST verify the human-readable labels for account, access key, secret key, session token, and expiration remain present

#### Scenario: Shell export output is stable
- **WHEN** `export` receives deterministic fake credentials
- **THEN** tests MUST verify the POSIX `export AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN` lines remain stable

#### Scenario: PowerShell export output is stable
- **WHEN** `export-ps` receives deterministic fake credentials
- **THEN** tests MUST verify the PowerShell `$env:` assignments remain stable

#### Scenario: Credential process output is stable
- **WHEN** `helper` receives deterministic fake credentials
- **THEN** tests MUST verify the credential_process-compatible JSON keys, version value, and RFC3339 expiration remain stable

### Requirement: Command behavior is characterized
The test suite SHALL preserve behavior for commands that affect processes, files, or AWS SSO list APIs while using fakes and temporary directories.

#### Scenario: Exec command builds expected environment
- **WHEN** `exec -- <command>` receives deterministic fake credentials and region
- **THEN** tests MUST verify the executed command would receive the AWS credential variables and `AWS_DEFAULT_REGION` without replacing the current test process

#### Scenario: Set command writes temporary credentials to fixture files
- **WHEN** `set PROFILE` receives deterministic fake credentials and a temporary AWS config directory
- **THEN** tests MUST verify the credentials and config files are updated with the expected sections and keys

#### Scenario: List accounts output is stable
- **WHEN** `list accounts` receives deterministic fake account data
- **THEN** tests MUST verify the tabular headers and account rows remain stable

#### Scenario: List roles output is stable
- **WHEN** `list roles ACCOUNT_ID` receives deterministic fake role data
- **THEN** tests MUST verify the tabular headers and role rows remain stable

### Requirement: Tests are isolated and safe by default
The test suite SHALL avoid external side effects and accidental credential exposure.

#### Scenario: Tests do not use real AWS state
- **WHEN** the unit test suite runs
- **THEN** tests MUST use temporary home directories and fake AWS clients instead of real AWS config, real SSO cache files, real credentials, or network calls

#### Scenario: Tests avoid leaking secrets
- **WHEN** fixtures contain fake credentials or access tokens
- **THEN** tests MUST use obvious non-secret placeholder values and MUST NOT require real secret material
