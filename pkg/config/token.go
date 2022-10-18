package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetSSOToken loops through all the caches files and extracts a valid token to use
func GetSSOToken(files []fs.DirEntry, ssoConfig SSOConfig, homedir string) (string, error) {


	if len(files) > 0 {
		// loop through all the files
		for _, file := range files {
			// read the contents into a JSON byte
			jsonContent, err := os.ReadFile(filepath.Join(homedir, ".aws", "sso", "cache", file.Name()))
			if err != nil {
				return "", fmt.Errorf("error reading aws SSO cache file: %v", err)
			}

			// initialize some SSOCacheConfig
			var cacheData SSOCacheConfig

			if err := json.Unmarshal(jsonContent, &cacheData); err != nil {
				return "", fmt.Errorf("error marshalling JSON data from cache file: %v", err)
			}

			// check if the file has a start url, if it doesn't, ignore it
			if cacheData.StartURL == ssoConfig.StartURL {
				// check if the file has an expiry time, if it doesn't ignore it
				if cacheData.ExpiresAt != "" {
					t, err := time.Parse(time.RFC3339, strings.Replace(cacheData.ExpiresAt, "UTC", "+00:00", -1))
					if err != nil {
						continue
					}
					if t.Unix() > time.Now().Unix() {
						return cacheData.AccessToken, nil
					}

				}
			}
			continue
		}
	}

	return "", fmt.Errorf("no valid cache files found, you might need to run aws sso login")

}
