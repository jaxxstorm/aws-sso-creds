## Context

`pkg/config.GetSSOToken` currently receives a directory listing and scans every file in `~/.aws/sso/cache` until it finds a valid token whose `startUrl` matches the profile's configured SSO start URL. AWS CLI cache filenames are often deterministic: the filename is the lowercase hex SHA1 digest of the exact SSO start URL, plus `.json`.

The implementation can therefore try the one expected cache file directly before doing any directory scan. Real-world AWS CLI SSO cache layouts can also contain usable tokens under other cache keys, so scanning must remain as a fallback to avoid breaking existing logged-in users.

## Goals / Non-Goals

**Goals:**

- Derive the cache filename as `hex(sha1([]byte(startURL))) + ".json"`.
- Read the derived cache file first for token lookup.
- Fall back to the existing cache scan when the derived cache file is absent or does not contain a valid usable token.
- Preserve validation for matching start URL, expiration parsing, expired token handling, malformed JSON, and login-oriented failures.
- Preserve user-facing credential command output and the existing login-oriented no-valid-cache-files error.
- Update tests to cover the deterministic hash example, direct lookup behavior, fallback behavior, and no-scanning behavior when the deterministic file is valid.

**Non-Goals:**

- Do not change AWS profile parsing semantics.
- Do not change credential output formats or command flags.
- Do not change AWS SDK calls or credential_process JSON.

## Decisions

### Decision: Add a cache filename helper in `pkg/config`

Introduce a small helper such as `SSOCacheFileName(startURL string) string` or an unexported equivalent. Tests should pin the example `https://d-xxxxxxxxxx.awsapps.com/start` to `5c26431228bc0d538e12104a3cbc37975e46c8f9.json`.

Alternative considered: inline SHA1 derivation inside token retrieval. A helper makes the behavior easy to test and keeps the retrieval code readable.

### Decision: Use deterministic lookup as a fast path, then scan for compatibility

Change token retrieval to construct `~/.aws/sso/cache/<derived>.json`, read it, unmarshal it, and run the same validation rules on that selected cache file. If that fast path does not return a usable token, scan `~/.aws/sso/cache` using the previous matching behavior so existing session-keyed cache files still work.

Alternative considered: remove scanning entirely. That is faster and simpler, but it breaks existing AWS CLI cache layouts where the usable token is not stored under the profile start URL hash.

### Decision: Preserve the existing public failure language

Missing file, start URL mismatch, expired token, and unparsable expiration should still produce the existing no-valid-cache-files message. A malformed selected cache file should still produce the existing cache JSON parsing context.

Alternative considered: introduce more precise errors such as "cache file missing." That may be nicer later, but this change should stay low-risk for scripts and users.

### Decision: Update callers to stop passing directory entries

The existing API accepts `[]fs.DirEntry`; callers no longer need to pass it. Update credential retrieval and list commands to call a start-URL-based lookup helper directly. The lookup helper owns both the deterministic fast path and fallback scan. Keep changes scoped to `pkg/config`, `pkg/credentials`, and list command token retrieval.

Alternative considered: keep the `[]fs.DirEntry` parameter but ignore it. That would reduce edits but leave misleading API shape around.

## Risks / Trade-offs

- Some AWS CLI cache layouts use a different cache key despite containing a matching `startUrl` -> Keep the old scan as fallback and test that behavior.
- Exact-string hashing means whitespace or URL case changes produce different filenames -> This matches the AWS CLI cache convention and avoids guessing normalization rules.
- Fallback scanning means malformed unrelated cache files can still affect lookup when the deterministic file is absent or invalid -> This preserves existing behavior for compatibility; direct lookup avoids that blast radius when the deterministic file is valid.
- Existing tests expect scanning behavior -> Update them to document deterministic lookup, fallback scanning, and the regression that unrelated malformed files are ignored when the direct file is valid.

## Migration Plan

Implement the helper and update token lookup, then adjust credential/list callers and tests. Run `go test ./...` and a small manual hash check against the documented example. No data migration is required.

## Open Questions

None.
