package credentials

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
)

type fakeRoleCredentialsClient struct {
	input *sso.GetRoleCredentialsInput
	out   *sso.GetRoleCredentialsOutput
	err   error
}

func (f *fakeRoleCredentialsClient) GetRoleCredentials(_ context.Context, in *sso.GetRoleCredentialsInput, _ ...func(*sso.Options)) (*sso.GetRoleCredentialsOutput, error) {
	f.input = in
	return f.out, f.err
}

func TestGetSSOCredentialsUsesParsedConfigTokenAndRegion(t *testing.T) {
	home := testutil.AWSHome(t)
	writeCredentialFixtures(t, home)

	client := &fakeRoleCredentialsClient{
		out: &sso.GetRoleCredentialsOutput{
			RoleCredentials: &types.RoleCredentials{
				AccessKeyId:     aws.String("fake-access-key-id"),
				SecretAccessKey: aws.String("fake-secret"),
				SessionToken:    aws.String("fake-session"),
				Expiration:      1893456000000,
			},
		},
	}
	var gotRegion string
	restoreCredentialSeams(t, client, &gotRegion, nil)

	creds, ssoConfig, awsCfg, err := GetSSOCredentials("dev", home)
	if err != nil {
		t.Fatalf("GetSSOCredentials returned error: %v", err)
	}

	if *creds.RoleCredentials.AccessKeyId != "fake-access-key-id" {
		t.Fatalf("unexpected access key %q", *creds.RoleCredentials.AccessKeyId)
	}
	if ssoConfig.AccountID != "123456789012" || ssoConfig.RoleName != "Admin" {
		t.Fatalf("unexpected SSO config %#v", ssoConfig)
	}
	if awsCfg.Region != "us-east-1" {
		t.Fatalf("unexpected AWS config region %q", awsCfg.Region)
	}
	if gotRegion != "us-west-2" {
		t.Fatalf("expected SSO client region us-west-2, got %q", gotRegion)
	}
	if *client.input.AccessToken != "fixture-access-token" {
		t.Fatalf("unexpected access token %q", *client.input.AccessToken)
	}
	if *client.input.AccountId != "123456789012" {
		t.Fatalf("unexpected account ID %q", *client.input.AccountId)
	}
	if *client.input.RoleName != "Admin" {
		t.Fatalf("unexpected role name %q", *client.input.RoleName)
	}
}

func TestGetSSOCredentialsMissingCacheDirectory(t *testing.T) {
	home := t.TempDir()
	testutil.WriteAWSConfig(t, home, `
[profile dev]
sso_start_url = https://example.awsapps.com/start
sso_region = us-west-2
sso_account_id = 123456789012
sso_role_name = Admin
`)

	restoreCredentialSeams(t, &fakeRoleCredentialsClient{}, nil, nil)

	_, _, _, err := GetSSOCredentials("dev", home)
	if err == nil || !strings.Contains(err.Error(), "error retrieving SSO token from cache files: no valid cache files found, you might need to run aws sso login") {
		t.Fatalf("expected missing token cache error, got %v", err)
	}
}

func TestGetSSOCredentialsWrapsAWSError(t *testing.T) {
	home := testutil.AWSHome(t)
	writeCredentialFixtures(t, home)

	restoreCredentialSeams(t, &fakeRoleCredentialsClient{err: errors.New("boom")}, nil, nil)

	_, _, _, err := GetSSOCredentials("dev", home)
	if err == nil || !strings.Contains(err.Error(), "error retrieving credentials from AWS: boom") {
		t.Fatalf("expected wrapped AWS error, got %v", err)
	}
}

func writeCredentialFixtures(t *testing.T, home string) {
	t.Helper()

	testutil.WriteAWSConfig(t, home, `
[profile dev]
sso_start_url = https://example.awsapps.com/start
sso_region = us-west-2
sso_account_id = 123456789012
sso_role_name = Admin
`)
	testutil.WriteSSOCache(t, home, cfg.SSOCacheFileName("https://example.awsapps.com/start"), `{
  "startUrl": "https://example.awsapps.com/start",
  "accessToken": "fixture-access-token",
  "expiresAt": "2999-01-02T03:04:05Z"
}`)
}

func restoreCredentialSeams(t *testing.T, client *fakeRoleCredentialsClient, gotRegion *string, loadErr error) {
	t.Helper()

	originalLoadDefaultConfig := loadDefaultConfig
	originalNewSSOClient := newSSOClient
	t.Cleanup(func() {
		loadDefaultConfig = originalLoadDefaultConfig
		newSSOClient = originalNewSSOClient
	})

	loadDefaultConfig = func(context.Context, ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
		if loadErr != nil {
			return aws.Config{}, loadErr
		}
		return aws.Config{Region: "us-east-1"}, nil
	}
	newSSOClient = func(_ aws.Config, ssoRegion string) getRoleCredentialsAPI {
		if gotRegion != nil {
			*gotRegion = ssoRegion
		}
		return client
	}
}
