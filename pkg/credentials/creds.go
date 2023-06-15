package credentials

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sso"
	"github.com/thuannfq/aws-sso-creds/pkg/config"
)

func GetSSOCredentials(profile string, homedir string) (*sso.GetRoleCredentialsOutput, string, error) {

	ssoConfig, err := config.GetSSOConfig(profile, homedir)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving SSO config: %w", err)
	}

	cacheFiles, err := os.ReadDir(filepath.Join(homedir, ".aws", "sso", "cache"))
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
	}

	token, err := config.GetSSOToken(cacheFiles, *ssoConfig, homedir)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving SSO token from cache files: %w", err)
	}

	sess := session.Must(session.NewSession())
	svc := sso.New(sess, aws.NewConfig().WithRegion(ssoConfig.Region))

	creds, err := svc.GetRoleCredentials(&sso.GetRoleCredentialsInput{
		AccessToken: &token,
		AccountId:   &ssoConfig.AccountID,
		RoleName:    &ssoConfig.RoleName,
	})

	if err != nil {
		return nil, "", fmt.Errorf("error retrieving credentials from AWS: %w", err)
	}

	return creds, ssoConfig.AccountID, nil

}
