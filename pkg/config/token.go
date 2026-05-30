package config

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
)

const noValidCacheFilesError = "no valid cache files found, you might need to run aws sso login"

var (
	errInvalidSSOCache    = errors.New(noValidCacheFilesError)
	refreshSSOAccessToken = refreshSSOAccessTokenWithOIDC
)

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
	if !errors.Is(err, errInvalidSSOCache) {
		return "", err
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
		if !cacheData.canRefresh() {
			return "", errInvalidSSOCache
		}
		return refreshSSOAccessToken(cacheData)
	}

	return cacheData.AccessToken, nil
}

func (c SSOCacheConfig) canRefresh() bool {
	return c.Region != "" && c.ClientID != "" && c.ClientSecret != "" && c.RefreshToken != ""
}

func refreshSSOAccessTokenWithOIDC(cacheData SSOCacheConfig) (string, error) {
	client := ssooidc.New(ssooidc.Options{Region: cacheData.Region})
	token, err := client.CreateToken(context.TODO(), &ssooidc.CreateTokenInput{
		ClientId:     aws.String(cacheData.ClientID),
		ClientSecret: aws.String(cacheData.ClientSecret),
		GrantType:    aws.String("refresh_token"),
		RefreshToken: aws.String(cacheData.RefreshToken),
	})
	if err != nil {
		return "", fmt.Errorf("error refreshing SSO access token: %w", err)
	}
	if token.AccessToken == nil || *token.AccessToken == "" {
		return "", errInvalidSSOCache
	}

	return *token.AccessToken, nil
}
