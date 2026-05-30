package accounts

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/jaxxstorm/aws-sso-creds/internal/testcreds"
	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/spf13/viper"
)

type fakeListAccountsClient struct {
	input *sso.ListAccountsInput
}

func (f *fakeListAccountsClient) ListAccounts(_ context.Context, in *sso.ListAccountsInput, _ ...func(*sso.Options)) (*sso.ListAccountsOutput, error) {
	f.input = in
	return &sso.ListAccountsOutput{
		AccountList: []types.AccountInfo{
			{
				AccountId:    aws.String("111122223333"),
				AccountName:  aws.String("dev-sandbox"),
				EmailAddress: aws.String("dev@example.com"),
			},
			{
				AccountId:    aws.String("444455556666"),
				AccountName:  aws.String("prod"),
				EmailAddress: aws.String("prod@example.com"),
			},
		},
	}, nil
}

func TestCommandOutput(t *testing.T) {
	home := testutil.AWSHome(t)
	writeAccountsFixtures(t, home)
	client := &fakeListAccountsClient{}
	restoreAccountsCommandState(t, home, client)

	output, err := testutil.CaptureStdout(t, Command().Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	for _, want := range []string{
		"ID",
		"NAME",
		"EMAIL ADDRESS",
		"111122223333",
		"dev-sandbox",
		"dev@example.com",
		"444455556666",
		"prod",
		"prod@example.com",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, output)
		}
	}
	if *client.input.AccessToken != "fixture-access-token" {
		t.Fatalf("unexpected access token %q", *client.input.AccessToken)
	}
	if *client.input.MaxResults != 10 {
		t.Fatalf("unexpected max results %d", *client.input.MaxResults)
	}
}

func writeAccountsFixtures(t *testing.T, home string) {
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

func restoreAccountsCommandState(t *testing.T, home string, client *fakeListAccountsClient) {
	t.Helper()

	originalLoadDefaultConfig := loadDefaultConfig
	originalNewSSOClient := newSSOClient
	t.Cleanup(func() {
		loadDefaultConfig = originalLoadDefaultConfig
		newSSOClient = originalNewSSOClient
		viper.Reset()
	})

	viper.Set("profile", "dev")
	viper.Set("home-directory", home)
	loadDefaultConfig = func(context.Context, ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
		return aws.Config{Region: testcreds.FakeRegion}, nil
	}
	newSSOClient = func(aws.Config) listAccountsAPI {
		return client
	}
}
