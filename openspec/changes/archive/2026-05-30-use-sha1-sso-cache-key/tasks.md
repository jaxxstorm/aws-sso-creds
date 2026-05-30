## 1. Token Lookup Tests

- [x] 1.1 Add a `pkg/config` test that proves `https://d-xxxxxxxxxx.awsapps.com/start` maps to `5c26431228bc0d538e12104a3cbc37975e46c8f9.json`.
- [x] 1.2 Update token selection tests to write cache fixtures under deterministic SHA1-derived filenames.
- [x] 1.3 Add `pkg/config` tests for missing deterministic cache file, malformed selected cache file, start URL mismatch, expired selected token, and unparsable selected expiration.
- [x] 1.4 Add a regression test proving unrelated malformed cache files are not read when the deterministic cache file is valid.
- [x] 1.5 Add regression tests proving token lookup falls back to matching non-deterministic cache files when the direct cache file is missing or stale.

## 2. Token Lookup Implementation

- [x] 2.1 Add a cache filename helper in `pkg/config` that hashes the exact SSO start URL bytes with SHA1 and appends `.json`.
- [x] 2.2 Replace caller-owned directory scanning with direct reading of `~/.aws/sso/cache/<sha1(start-url)>.json` as the first lookup path.
- [x] 2.3 Preserve existing selected-file JSON parsing, start URL validation, expiration parsing, expired-token handling, and no-valid-cache-files error text.
- [x] 2.4 Preserve cache directory scanning as an internal compatibility fallback when the deterministic file cannot provide a valid token.

## 3. Caller Updates

- [x] 3.1 Update `pkg/credentials` so credential retrieval no longer reads or passes cache directory entries.
- [x] 3.2 Update `list accounts` token retrieval to use deterministic cache lookup without directory scanning.
- [x] 3.3 Update `list roles` token retrieval to use deterministic cache lookup without directory scanning.
- [x] 3.4 Update credential/list tests affected by the token lookup API shape.

## 4. Verification

- [x] 4.1 Run `go test ./...` and fix failures without changing user-visible command output.
- [x] 4.2 Manually verify the documented SHA1 example with `printf %s "https://d-xxxxxxxxxx.awsapps.com/start" | shasum -a 1`.
- [x] 4.3 Review changes to confirm callers no longer scan the cache directory and deterministic lookup is tried before fallback scanning.
