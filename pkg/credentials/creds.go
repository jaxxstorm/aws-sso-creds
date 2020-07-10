package credentials

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/pkg/config"
)

func GetSSOCredentials(profile string, homedir string) (*sso.GetRoleCredentialsOutput, error) {

	ssoConfig, err := config.GetSSOConfig(profile, homedir)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SSO config: %w", err)
	}

	cacheFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/.aws/sso/cache", homedir))
	if err != nil {
		return nil, fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
	}

	token, err := config.GetSSOToken(cacheFiles, *ssoConfig, homedir)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SSO token from cache files: %w", err)
	}

	sess := session.Must(session.NewSession())
	svc := sso.New(sess, aws.NewConfig().WithRegion(ssoConfig.Region))

	creds, err := svc.GetRoleCredentials(&sso.GetRoleCredentialsInput{
		AccessToken: &token,
		AccountId:   &ssoConfig.AccountID,
		RoleName:    &ssoConfig.RoleName,
	})

	if err != nil {
		return nil, fmt.Errorf("error retrieving credentials from AWS: %w", err)
	}

	return creds, nil

}
