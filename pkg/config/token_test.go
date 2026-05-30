package config

import (
	"strings"
	"testing"

	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
)

const testStartURL = "https://example.awsapps.com/start"

func TestSSOCacheFileName(t *testing.T) {
	got := SSOCacheFileName("https://d-xxxxxxxxxx.awsapps.com/start")
	want := "5c26431228bc0d538e12104a3cbc37975e46c8f9.json"
	if got != want {
		t.Fatalf("unexpected cache file name %q, want %q", got, want)
	}
}

func TestSSOCacheFileNameUsesExactURLBytes(t *testing.T) {
	original := SSOCacheFileName("https://example.awsapps.com/start")
	changedCase := SSOCacheFileName("https://EXAMPLE.awsapps.com/start")
	withSpace := SSOCacheFileName("https://example.awsapps.com/start ")

	if original == changedCase {
		t.Fatalf("expected case-changed URL to hash differently")
	}
	if original == withSpace {
		t.Fatalf("expected whitespace-changed URL to hash differently")
	}
}

func TestGetSSOTokenSelectsDeterministicMatchingToken(t *testing.T) {
	home := testutil.AWSHome(t)
	writeDeterministicCache(t, home, testStartURL, `{
  "startUrl": "https://example.awsapps.com/start",
  "region": "us-west-2",
  "accessToken": "fixture-access-token",
  "expiresAt": "2999-01-02T03:04:05Z"
}`)

	token, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	if err != nil {
		t.Fatalf("GetSSOToken returned error: %v", err)
	}
	if token != "fixture-access-token" {
		t.Fatalf("unexpected token %q", token)
	}
}

func TestGetSSOTokenMissingDeterministicCacheFile(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteSSOCache(t, home, "unrelated.json", validCacheJSON("https://other.awsapps.com/start", "other-token", "2999-01-02T03:04:05Z"))

	_, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	assertNoValidCacheFilesError(t, err)
}

func TestGetSSOTokenFallsBackToMatchingCacheFileWhenDeterministicMissing(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteSSOCache(t, home, "session-keyed-cache.json", validCacheJSON(testStartURL, "fallback-token", "2999-01-02T03:04:05Z"))

	token, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	if err != nil {
		t.Fatalf("GetSSOToken returned error: %v", err)
	}
	if token != "fallback-token" {
		t.Fatalf("unexpected token %q", token)
	}
}

func TestGetSSOTokenReturnsMalformedSelectedJSONError(t *testing.T) {
	home := testutil.AWSHome(t)
	writeDeterministicCache(t, home, testStartURL, `{not-json`)

	_, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	if err == nil || !strings.Contains(err.Error(), "error marshalling JSON data from cache file") {
		t.Fatalf("expected malformed JSON error, got %v", err)
	}
}

func TestGetSSOTokenRejectsSelectedStartURLMismatch(t *testing.T) {
	home := testutil.AWSHome(t)
	writeDeterministicCache(t, home, testStartURL, validCacheJSON("https://other.awsapps.com/start", "wrong-token", "2999-01-02T03:04:05Z"))

	_, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	assertNoValidCacheFilesError(t, err)
}

func TestGetSSOTokenRejectsExpiredSelectedToken(t *testing.T) {
	home := testutil.AWSHome(t)
	writeDeterministicCache(t, home, testStartURL, validCacheJSON(testStartURL, "expired-token", "2000-01-02T03:04:05Z"))

	_, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	assertNoValidCacheFilesError(t, err)
}

func TestGetSSOTokenFallsBackWhenDeterministicTokenExpired(t *testing.T) {
	home := testutil.AWSHome(t)
	writeDeterministicCache(t, home, testStartURL, validCacheJSON(testStartURL, "expired-token", "2000-01-02T03:04:05Z"))
	testutil.WriteSSOCache(t, home, "session-keyed-cache.json", validCacheJSON(testStartURL, "fresh-fallback-token", "2999-01-02T03:04:05Z"))

	token, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	if err != nil {
		t.Fatalf("GetSSOToken returned error: %v", err)
	}
	if token != "fresh-fallback-token" {
		t.Fatalf("unexpected token %q", token)
	}
}

func TestGetSSOTokenRejectsUnparseableSelectedExpiration(t *testing.T) {
	home := testutil.AWSHome(t)
	writeDeterministicCache(t, home, testStartURL, validCacheJSON(testStartURL, "bad-expiry-token", "not-a-time"))

	_, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	assertNoValidCacheFilesError(t, err)
}

func TestGetSSOTokenIgnoresMalformedUnrelatedCacheFiles(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteSSOCache(t, home, "malformed-unrelated.json", `{not-json`)
	writeDeterministicCache(t, home, testStartURL, validCacheJSON(testStartURL, "fixture-access-token", "2999-01-02T03:04:05UTC"))

	token, err := GetSSOToken(SSOConfig{StartURL: testStartURL}, home)
	if err != nil {
		t.Fatalf("GetSSOToken returned error: %v", err)
	}
	if token != "fixture-access-token" {
		t.Fatalf("unexpected token %q", token)
	}
}

func writeDeterministicCache(t *testing.T, home string, startURL string, content string) {
	t.Helper()
	testutil.WriteSSOCache(t, home, SSOCacheFileName(startURL), content)
}

func validCacheJSON(startURL string, token string, expiresAt string) string {
	return `{
  "startUrl": "` + startURL + `",
  "accessToken": "` + token + `",
  "expiresAt": "` + expiresAt + `"
}`
}

func assertNoValidCacheFilesError(t *testing.T, err error) {
	t.Helper()

	if err == nil || err.Error() != noValidCacheFilesError {
		t.Fatalf("expected no valid cache files error, got %v", err)
	}
}
