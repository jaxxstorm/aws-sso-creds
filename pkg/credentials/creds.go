package credentials

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
)

func GetSSOCredentials(profile string, homedir string) (*sso.GetRoleCredentialsOutput, string, error) {
	ssoConfig, err := cfg.GetSSOConfig(profile, homedir)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving SSO config: %w", err)
	}

	cacheFiles, err := os.ReadDir(filepath.Join(homedir, ".aws", "sso", "cache"))
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
	}

	token, err := cfg.GetSSOToken(cacheFiles, *ssoConfig, homedir)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving SSO token from cache files: %w", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(ssoConfig.Region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, "", fmt.Errorf("error loading AWS configuration: %w", err)
	}

	svc := sso.NewFromConfig(cfg)

	creds, err := svc.GetRoleCredentials(context.TODO(), &sso.GetRoleCredentialsInput{
		AccessToken: aws.String(token),
		AccountId:   aws.String(ssoConfig.AccountID),
		RoleName:    aws.String(ssoConfig.RoleName),
	})
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving credentials from AWS: %w", err)
	}

	return creds, ssoConfig.AccountID, nil
}
