package roles

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

type fakeListAccountRolesClient struct {
	input *sso.ListAccountRolesInput
}

func (f *fakeListAccountRolesClient) ListAccountRoles(_ context.Context, in *sso.ListAccountRolesInput, _ ...func(*sso.Options)) (*sso.ListAccountRolesOutput, error) {
	f.input = in
	return &sso.ListAccountRolesOutput{
		RoleList: []types.RoleInfo{
			{RoleName: aws.String("AdministratorAccess")},
			{RoleName: aws.String("ReadOnly")},
		},
	}, nil
}

func TestCommandOutput(t *testing.T) {
	home := testutil.AWSHome(t)
	writeRolesFixtures(t, home)
	client := &fakeListAccountRolesClient{}
	restoreRolesCommandState(t, home, client)

	command := Command()
	command.SetArgs([]string{"111122223333"})
	output, err := testutil.CaptureStdout(t, command.Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	for _, want := range []string{
		"ROLE NAME",
		"AdministratorAccess",
		"ReadOnly",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, output)
		}
	}
	if *client.input.AccessToken != "fixture-access-token" {
		t.Fatalf("unexpected access token %q", *client.input.AccessToken)
	}
	if *client.input.AccountId != "111122223333" {
		t.Fatalf("unexpected account ID %q", *client.input.AccountId)
	}
	if *client.input.MaxResults != 10 {
		t.Fatalf("unexpected max results %d", *client.input.MaxResults)
	}
}

func writeRolesFixtures(t *testing.T, home string) {
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

func restoreRolesCommandState(t *testing.T, home string, client *fakeListAccountRolesClient) {
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
	newSSOClient = func(aws.Config) listAccountRolesAPI {
		return client
	}
}
