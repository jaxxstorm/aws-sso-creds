## Why

AWS SSO cache filenames are commonly deterministic SHA1 hex digests of SSO start URLs, so `aws-sso-creds` can usually find the token for a profile without scanning every file under `~/.aws/sso/cache`. Looking up the exact expected cache file first is simpler, faster, and avoids being affected by unrelated or malformed cache files when the deterministic cache file is present.

Some existing AWS CLI SSO cache layouts store usable tokens under other cache keys, especially with `sso_session`-based profiles. To preserve the workflow that has worked for years, the direct lookup should be a fast path with the previous scanning behavior retained as a compatibility fallback.

## What Changes

- Compute the expected AWS SSO cache filename as the lowercase hex-encoded SHA1 digest of the configured `sso_start_url`, plus `.json`.
- Try `~/.aws/sso/cache/<sha1(start-url)>.json` first when retrieving an access token for a profile.
- Fall back to the existing cache directory scan when the deterministic file cannot provide a valid token.
- Preserve existing token validation behavior, including start URL matching, expiration parsing, expired-token handling, malformed JSON handling during scanning, and the login-oriented no-valid-cache error.
- Update credential retrieval so callers use a single token lookup API instead of listing and scanning the cache directory themselves.
- Update unit tests to cover SHA1 filename derivation, direct file lookup, compatibility fallback lookup, missing cache behavior, malformed selected cache file behavior, and ignoring unrelated malformed cache files when the deterministic cache file is valid.
- Non-goals: change AWS profile parsing, alter credential output formats, change AWS authentication semantics, or remove compatibility with existing non-SHA1 cache keys.

## Capabilities

### New Capabilities

- `deterministic-sso-cache-lookup`: Resolve AWS SSO cache tokens by deriving the expected cache file path from the profile start URL before falling back to compatible cache scanning.

### Modified Capabilities

None.

## Impact

- Affected code: `pkg/config/token.go`, `pkg/credentials/creds.go`, and tests around token selection and credential retrieval.
- Affected CLI commands: every command that retrieves AWS SSO credentials (`get`, `export`, `export-ps`, `exec`, `helper`, `set`, `list accounts`, and `list roles`) benefits from the direct cache lookup when available; user-visible output should remain unchanged.
- Affected behavior: unrelated cache files should no longer influence token retrieval when the deterministic cache file is valid, while existing cache layouts that require scanning should keep working.
- Dependencies: use Go standard library SHA1/hex/path handling; no new external dependency is needed.
