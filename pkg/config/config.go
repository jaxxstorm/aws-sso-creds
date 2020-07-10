package config

import (
	"fmt"
	"github.com/bigkevmcd/go-configparser"
)


// GetSSOConfig retrieves the SSO configuration for a given AWS profile
func GetSSOConfig(profile string, homedir string) (*SSOConfig, error) {

	// parse the configuration file
	p, err := configparser.NewConfigParserFromFile(fmt.Sprintf("%s/.aws/config", homedir))

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
