package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// GetSSOToken loops through all the caches files and extracts a valid token to use
func GetSSOToken(files []os.FileInfo, ssoConfig SSOConfig, homedir string) (string, error) {

	if len(files) > 0 {
		// loop through all the files
		for _, file := range files {
			// read the contents into a JSON byte
			jsonContent, err := ioutil.ReadFile(fmt.Sprintf("%s/.aws/sso/cache/%s", homedir, file.Name()))
			if err != nil {
				panic(err)
			}

			// initialize some SSOCacheConfig
			var cacheData SSOCacheConfig

			json.Unmarshal(jsonContent, &cacheData)

			// check if the file has a start url, if it doesn't, ignore it
			if cacheData.StartUrl == ssoConfig.StartUrl {
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
