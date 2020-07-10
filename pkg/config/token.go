package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// GetSSOToken loops through all the caches files and extracts a valid token to use
func GetSSOToken(files []os.FileInfo, ssoConfig SSOConfig) (string, error) {

	if len(files) > 0 {
		// loop through all the files
		for _, file := range files {
			// read the contents into a JSON byte
			jsonContent, err := ioutil.ReadFile(fmt.Sprintf("/Users/lbriggs/.aws/sso/cache/%s", file.Name()))
			if err != nil {
				panic(err)
			}

			// initialize some SSOCacheConfig
			var cacheData SSOCacheConfig

			json.Unmarshal(jsonContent, &cacheData)

			if (cacheData.StartUrl == ssoConfig.StartUrl) {
				return cacheData.AccessToken, nil
			}

			return "", fmt.Errorf("Unable to find a valid access token")

		}
	}

	return "", fmt.Errorf("No cache files found")

}
