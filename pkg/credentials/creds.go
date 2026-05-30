package credentials

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
)

type getRoleCredentialsAPI interface {
	GetRoleCredentials(context.Context, *sso.GetRoleCredentialsInput, ...func(*sso.Options)) (*sso.GetRoleCredentialsOutput, error)
}

var (
	loadDefaultConfig = awsconfig.LoadDefaultConfig
	newSSOClient      = func(awsCfg aws.Config, ssoRegion string) getRoleCredentialsAPI {
		return sso.NewFromConfig(awsCfg, func(o *sso.Options) {
			// We specify the SSO region here because it applies only to this client instance.
			o.Region = ssoRegion
		})
	}
)

func GetSSOCredentials(profile string, homedir string) (*sso.GetRoleCredentialsOutput, *cfg.SSOConfig, *aws.Config, error) {
	ssoConfig, err := cfg.GetSSOConfig(profile, homedir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving SSO config: %w", err)
	}

	token, err := cfg.GetSSOToken(*ssoConfig, homedir)
	if err != nil {
		return nil, ssoConfig, nil, fmt.Errorf("error retrieving SSO token from cache files: %w", err)
	}

	awsCfg, err := loadDefaultConfig(context.TODO(),
		// Don't set the region here in config because it would affect all clients created from this config
		awsconfig.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, ssoConfig, nil, fmt.Errorf("error loading AWS configuration: %w", err)
	}

	svc := newSSOClient(awsCfg, ssoConfig.Region)

	creds, err := svc.GetRoleCredentials(context.TODO(), &sso.GetRoleCredentialsInput{
		AccessToken: aws.String(token),
		AccountId:   aws.String(ssoConfig.AccountID),
		RoleName:    aws.String(ssoConfig.RoleName),
	})
	if err != nil {
		return nil, ssoConfig, &awsCfg, fmt.Errorf("error retrieving credentials from AWS: %w", err)
	}

	return creds, ssoConfig, &awsCfg, nil
}
