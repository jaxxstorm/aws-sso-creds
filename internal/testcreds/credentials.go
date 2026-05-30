package testcreds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
)

const (
	FakeAccessKeyID     = "fake-access-key-id"
	FakeSecretAccessKey = "fake-secret-access-key"
	FakeSessionToken    = "fake-session-token"
	FakeAccountID       = "123456789012"
	FakeRoleName        = "Admin"
	FakeRegion          = "us-west-2"
	FakeExpiration      = int64(1893456000000)
)

func FakeCredentialResult() (*sso.GetRoleCredentialsOutput, *cfg.SSOConfig, *aws.Config) {
	return &sso.GetRoleCredentialsOutput{
			RoleCredentials: &types.RoleCredentials{
				AccessKeyId:     aws.String(FakeAccessKeyID),
				SecretAccessKey: aws.String(FakeSecretAccessKey),
				SessionToken:    aws.String(FakeSessionToken),
				Expiration:      FakeExpiration,
			},
		},
		&cfg.SSOConfig{
			StartURL:  "https://example.awsapps.com/start",
			Region:    FakeRegion,
			AccountID: FakeAccountID,
			RoleName:  FakeRoleName,
		},
		&aws.Config{Region: FakeRegion}
}
