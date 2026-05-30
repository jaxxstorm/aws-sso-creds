package config

import (
	"strings"
	"testing"

	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
)

func TestGetSSOTokenSelectsValidMatchingToken(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteSSOCache(t, home, "token.json", `{
  "startUrl": "https://example.awsapps.com/start",
  "region": "us-west-2",
  "accessToken": "fixture-access-token",
  "expiresAt": "2999-01-02T03:04:05Z"
}`)

	token, err := GetSSOToken(readCacheEntries(t, home), SSOConfig{
		StartURL: "https://example.awsapps.com/start",
	}, home)
	if err != nil {
		t.Fatalf("GetSSOToken returned error: %v", err)
	}
	if token != "fixture-access-token" {
		t.Fatalf("unexpected token %q", token)
	}
}

func TestGetSSOTokenIgnoresExpiredAndNonMatchingTokens(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteSSOCache(t, home, "expired.json", `{
  "startUrl": "https://example.awsapps.com/start",
  "accessToken": "expired-token",
  "expiresAt": "2000-01-02T03:04:05Z"
}`)
	testutil.WriteSSOCache(t, home, "other-start-url.json", `{
  "startUrl": "https://other.awsapps.com/start",
  "accessToken": "other-token",
  "expiresAt": "2999-01-02T03:04:05Z"
}`)

	_, err := GetSSOToken(readCacheEntries(t, home), SSOConfig{
		StartURL: "https://example.awsapps.com/start",
	}, home)
	if err == nil || err.Error() != "no valid cache files found, you might need to run aws sso login" {
		t.Fatalf("expected no valid cache files error, got %v", err)
	}
}

func TestGetSSOTokenReturnsMalformedJSONError(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteSSOCache(t, home, "bad.json", `{not-json`)

	_, err := GetSSOToken(readCacheEntries(t, home), SSOConfig{
		StartURL: "https://example.awsapps.com/start",
	}, home)
	if err == nil || !strings.Contains(err.Error(), "error marshalling JSON data from cache file") {
		t.Fatalf("expected malformed JSON error, got %v", err)
	}
}

func TestGetSSOTokenSkipsUnparseableExpiration(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteSSOCache(t, home, "bad-expiry.json", `{
  "startUrl": "https://example.awsapps.com/start",
  "accessToken": "bad-expiry-token",
  "expiresAt": "not-a-time"
}`)
	testutil.WriteSSOCache(t, home, "valid.json", `{
  "startUrl": "https://example.awsapps.com/start",
  "accessToken": "later-valid-token",
  "expiresAt": "2999-01-02T03:04:05UTC"
}`)

	token, err := GetSSOToken(readCacheEntries(t, home), SSOConfig{
		StartURL: "https://example.awsapps.com/start",
	}, home)
	if err != nil {
		t.Fatalf("GetSSOToken returned error: %v", err)
	}
	if token != "later-valid-token" {
		t.Fatalf("unexpected token %q", token)
	}
}
