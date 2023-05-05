package config

import (
	"fmt"
	"path/filepath"

	"github.com/bigkevmcd/go-configparser"
)

// GetSSOConfig retrieves the SSO configuration for a given AWS profile
func GetSSOConfig(profile string, homedir string) (*SSOConfig, error) {

	// parse the configuration file
	p, err := configparser.NewConfigParserFromFile(filepath.Join(homedir, ".aws", "config"))

	if err != nil {
		return nil, err
	}

	// build a section name
	var section string
	if profile == "" || profile == "default" {
		section = "default"
		profile = "<default>"
	} else {
		section = fmt.Sprintf("profile %s", profile)
	}

	// FIXME: make this better
	if p.HasSection(section) {
		var startURL string
		var region string
		var roleName string

		// check if we have an sso_session section
		// if we do, retrieve the vars from that
		// if not, retrieve it from the profile
		ssoSession, err := p.Get(section, "sso_session")
		if err == nil {
			ssoSection := fmt.Sprintf("sso-session %s", ssoSession)

			startURL, err = p.Get(ssoSection, "sso_start_url")
			if err != nil {
				return nil, fmt.Errorf("no SSO url in sso-session: %s", ssoSection)
			}
			region, err = p.Get(ssoSection, "sso_region")
			if err != nil {
				return nil, fmt.Errorf("no SSO region in sso-session: %s", ssoSection)
			}
			roleName, err = p.Get(ssoSection, "sso_role_name")
			if err != nil {
				return nil, fmt.Errorf("no SSO role name in sso-session: %s", ssoSection)
			}
		} else {
			startURL, err = p.Get(section, "sso_start_url")
			if err != nil {
				return nil, fmt.Errorf("no SSO url in profile: %s", profile)
			}
			region, err = p.Get(section, "sso_region")
			if err != nil {
				return nil, fmt.Errorf("no SSO region in profile: %s", profile)
			}
			roleName, err = p.Get(section, "sso_role_name")
			if err != nil {
				return nil, fmt.Errorf("no SSO role name in profile: %s", profile)
			}
		}
		// account id is always going to exist within the profile
		ssoAccountID, err := p.Get(section, "sso_account_id")
		if err != nil {
			return nil, fmt.Errorf("no SSO account id in profile: %s", profile)
		}

		return &SSOConfig{
			StartURL:  startURL,
			Region:    region,
			AccountID: ssoAccountID,
			RoleName:  roleName,
		}, nil

	}

	return nil, fmt.Errorf("unable to find profile: %s", profile)

}
