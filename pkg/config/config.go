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

	if !p.HasSection(section) {
		return nil, fmt.Errorf("unable to find profile: %s", profile)
	}

	c := &SSOConfig{}

	// Check if we have a sso-session section and merge SSO options from the sso-session and profile
	// SSO options from the profile take precedence over the shared sso-session
	// Any SSO option can be present in either sso-session or profile
	ssoSession, err := p.Get(section, "sso_session")
	if err == nil {
		ssoSection := fmt.Sprintf("sso-session %s", ssoSession)
		mergeSSOConfig(p, ssoSection, c)

		// sso-session requires sso_start_url and sso_region
		if c.Region == "" {
			return nil, fmt.Errorf("no sso_region in sso-session %q", ssoSession)
		}
		if c.StartURL == "" {
			return nil, fmt.Errorf("no sso_start_url in sso-session %q", ssoSession)
		}
	}

	mergeSSOConfig(p, section, c)

	// Validate the required SSO options
	if c.Region == "" {
		return nil, fmt.Errorf("no sso_region in profile %q and its sso_session", profile)
	}
	if c.StartURL == "" {
		return nil, fmt.Errorf("no sso_start_url in profile %q and its sso_session", profile)
	}
	if c.AccountID == "" {
		return nil, fmt.Errorf("no sso_account_id in profile %q or its sso_session", profile)
	}
	if c.RoleName == "" {
		return nil, fmt.Errorf("no sso_role_name in profile %q or its sso_session", profile)
	}

	return c, nil
}

// mergeSSOConfig merges non-empty SSO options from the specified section (sso-session or profile) into the SSOConfig struct s overwriting the existing values
//
// TODO: Should be removed in favor of github.com/aws/aws-sdk-go-v2/config and github.com/aws/aws-sdk-go-v2/credentials
func mergeSSOConfig(p *configparser.ConfigParser, section string, s *SSOConfig) {
	if accountID, err := p.Get(section, "sso_account_id"); err == nil && accountID != "" {
		s.AccountID = accountID
	}

	if startURL, err := p.Get(section, "sso_start_url"); err == nil && startURL != "" {
		s.StartURL = startURL
	}

	if region, err := p.Get(section, "sso_region"); err == nil && region != "" {
		s.Region = region
	}

	if roleName, err := p.Get(section, "sso_role_name"); err == nil && roleName != "" {
		s.RoleName = roleName
	}
}
