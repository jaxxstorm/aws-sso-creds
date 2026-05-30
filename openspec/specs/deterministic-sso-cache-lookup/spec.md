## Purpose

Document how AWS SSO cache files are resolved from profile start URLs while preserving compatibility with existing AWS CLI cache layouts.

## Requirements

### Requirement: SSO cache filename is derived from start URL
The system SHALL derive the AWS SSO cache filename by computing the SHA1 digest of the exact configured SSO start URL bytes and formatting that digest as lowercase hexadecimal with a `.json` suffix.

#### Scenario: Start URL maps to expected SHA1 filename
- **WHEN** the configured SSO start URL is `https://d-xxxxxxxxxx.awsapps.com/start`
- **THEN** the expected cache filename MUST be `5c26431228bc0d538e12104a3cbc37975e46c8f9.json`

#### Scenario: URL bytes are not normalized before hashing
- **WHEN** a configured SSO start URL is used to derive the cache filename
- **THEN** the system MUST hash the exact configured start URL string rather than lowercasing, trimming, escaping, or otherwise normalizing it

### Requirement: Token lookup tries the deterministic cache file first
The system SHALL try the cache file derived from the configured SSO start URL before scanning other cache files when retrieving an AWS SSO access token for a profile.

#### Scenario: Matching deterministic cache file is used
- **WHEN** `~/.aws/sso/cache/<sha1(start-url)>.json` exists and contains an unexpired token for the configured start URL
- **THEN** the system MUST return that access token

#### Scenario: Missing deterministic cache file falls back to compatible scan
- **WHEN** the deterministic cache file for the configured start URL does not exist
- **AND** another cache file contains an unexpired token for the configured start URL
- **THEN** the system MUST return that access token

#### Scenario: Invalid deterministic cache file falls back to compatible scan
- **WHEN** the deterministic cache file exists but cannot provide a valid token for the configured start URL
- **AND** another cache file contains an unexpired token for the configured start URL
- **THEN** the system MUST return the fallback access token

#### Scenario: No matching cache file reports login guidance
- **WHEN** the deterministic cache file for the configured start URL does not exist
- **AND** no other cache file contains an unexpired token for the configured start URL
- **THEN** the system MUST return the existing no-valid-cache-files error that tells the user they might need to run `aws sso login`

#### Scenario: Malformed unrelated cache file is ignored
- **WHEN** an unrelated cache file contains malformed JSON but the deterministic cache file contains a valid token
- **THEN** token retrieval MUST succeed using the deterministic cache file

### Requirement: Cache file validation is preserved
The system SHALL preserve existing validation for cache files used during deterministic lookup and fallback scanning.

#### Scenario: Cache file start URL mismatch is rejected
- **WHEN** the deterministic cache file exists but its `startUrl` value does not match the configured SSO start URL
- **THEN** the system MUST reject the file and return the existing no-valid-cache-files error

#### Scenario: Cache file expired token is rejected
- **WHEN** the deterministic cache file exists but its expiration is in the past
- **THEN** the system MUST reject the token and return the existing no-valid-cache-files error

#### Scenario: Cache file malformed JSON returns parsing error
- **WHEN** the deterministic cache file exists but contains malformed JSON
- **THEN** the system MUST return the existing JSON parsing error context for cache files

#### Scenario: Cache file unparsable expiration is rejected
- **WHEN** the deterministic cache file exists but its expiration cannot be parsed
- **THEN** the system MUST reject the file and return the existing no-valid-cache-files error

### Requirement: Credential retrieval uses centralized token lookup
Credential retrieval SHALL use centralized SSO cache lookup instead of listing cache directory entries before calling AWS SSO APIs.

#### Scenario: Credential retrieval succeeds from deterministic cache file
- **WHEN** the profile config resolves to an SSO start URL and the corresponding deterministic cache file contains a valid token
- **THEN** credential retrieval MUST use that token to request role credentials from AWS SSO

#### Scenario: Credential retrieval succeeds from fallback cache file
- **WHEN** the profile config resolves to an SSO start URL and a non-deterministic cache file contains a valid token for that start URL
- **THEN** credential retrieval MUST use that token to request role credentials from AWS SSO

#### Scenario: Credential retrieval does not require listing cache directory entries
- **WHEN** retrieving credentials for a profile
- **THEN** callers outside `pkg/config` MUST NOT depend on iterating all entries in `~/.aws/sso/cache`
