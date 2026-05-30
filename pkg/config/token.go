package config

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const noValidCacheFilesError = "no valid cache files found, you might need to run aws sso login"

var errInvalidSSOCache = errors.New(noValidCacheFilesError)

// SSOCacheFileName returns the AWS SSO cache file name for a start URL.
func SSOCacheFileName(startURL string) string {
	sum := sha1.Sum([]byte(startURL))
	return hex.EncodeToString(sum[:]) + ".json"
}

// GetSSOToken reads the deterministic SSO cache file and extracts a valid token to use.
// If AWS has stored the token under another cache key, it falls back to scanning
// the cache directory for compatibility with existing AWS CLI cache layouts.
func GetSSOToken(ssoConfig SSOConfig, homedir string) (string, error) {
	cacheDir := filepath.Join(homedir, ".aws", "sso", "cache")
	token, err := getSSOTokenFromCacheFile(filepath.Join(cacheDir, SSOCacheFileName(ssoConfig.StartURL)), ssoConfig)
	if err == nil {
		return token, nil
	}

	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return "", errInvalidSSOCache
	}

	for _, file := range files {
		token, err := getSSOTokenFromCacheFile(filepath.Join(cacheDir, file.Name()), ssoConfig)
		if err == nil {
			return token, nil
		}
		if !errors.Is(err, errInvalidSSOCache) {
			return "", err
		}
	}

	return "", errInvalidSSOCache
}

func getSSOTokenFromCacheFile(path string, ssoConfig SSOConfig) (string, error) {
	jsonContent, err := os.ReadFile(path)
	if err != nil {
		return "", errInvalidSSOCache
	}

	var cacheData SSOCacheConfig
	if err := json.Unmarshal(jsonContent, &cacheData); err != nil {
		return "", fmt.Errorf("error marshalling JSON data from cache file: %v", err)
	}

	if cacheData.StartURL != ssoConfig.StartURL || cacheData.ExpiresAt == "" {
		return "", errInvalidSSOCache
	}

	t, err := time.Parse(time.RFC3339, strings.ReplaceAll(cacheData.ExpiresAt, "UTC", "+00:00"))
	if err != nil {
		return "", errInvalidSSOCache
	}
	if t.Unix() <= time.Now().Unix() {
		return "", errInvalidSSOCache
	}

	return cacheData.AccessToken, nil
}
