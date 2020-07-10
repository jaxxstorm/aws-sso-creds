package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"os"

	"github.com/bigkevmcd/go-configparser"
	"github.com/aws/aws-sdk-go/service/sso"
)

type SSOConfig struct {
	StartUrl  string
	Region    string
	AccountID string
	RoleName  string
}

type SSOCacheConfig struct {
	StartUrl    string `json:"startUrl"`
	Region      string `json:"region"`
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"` // FIXME not a string
}

func main() {

	profile := "pulumi-dev-sandbox"

	ssoConfig, err := getSSOConfig(profile)
	if err != nil {
		panic(err)
	}

	cacheFiles, err := ioutil.ReadDir("/Users/lbriggs/.aws/sso/cache")
	if err != nil {
		panic(err)
	}

	token, err := getSSOToken(cacheFiles, *ssoConfig)
	if err != nil {
		panic(err)
	}

	sess := session.Must(session.NewSession())
	svc := sso.New(sess, aws.NewConfig().WithRegion(ssoConfig.Region))

	creds, err := svc.GetRoleCredentials(&sso.GetRoleCredentialsInput{
		AccessToken: &token,
		AccountId: &ssoConfig.AccountID,
		RoleName: &ssoConfig.RoleName,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(creds)


}

// getSSOToken loops through all the caches files and extracts a valid token to use
func getSSOToken(files []os.FileInfo, ssoConfig SSOConfig) (string, error) {

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


// getSSOConfig retrieves the SSO configuration for a given AWS profile
func getSSOConfig(profile string) (*SSOConfig, error) {

	// parse the configuration file
	// FIXME: make this configurable or use go-homedir
	p, err := configparser.NewConfigParserFromFile("/Users/lbriggs/.aws/config")

	if err != nil {
		return nil, err
	}

	// build a section name
	section := fmt.Sprintf("profile %s", profile)

	// FIXME: make this better
	if p.HasSection(section) {
		ssoStartUrl, err := p.Get(section, "sso_start_url")
		if err != nil {
			fmt.Println("no SSO url in profile")
		}
		ssoRegion, err := p.Get(section, "sso_region")
		if err != nil {
			fmt.Println("no SSO region in profile")
		}
		ssoAccountId, err := p.Get(section, "sso_account_id")
		if err != nil {
			fmt.Println("no SSO account id in profile")
		}
		ssoRoleName, err := p.Get(section, "sso_role_name")
		if err != nil {
			fmt.Println("no SSO role name in profile")
		}

		return &SSOConfig{
			StartUrl:  ssoStartUrl,
			Region:    ssoRegion,
			AccountID: ssoAccountId,
			RoleName:  ssoRoleName,
		}, nil

	} else {
		return nil, fmt.Errorf("unable to find profile %s", profile)
	}
}
