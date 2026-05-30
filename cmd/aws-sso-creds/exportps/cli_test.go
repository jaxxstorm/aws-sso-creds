package exportps

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/internal/testcreds"
	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/spf13/viper"
)

func TestCommandOutput(t *testing.T) {
	restoreExportPSCommandState(t)

	output, err := testutil.CaptureStdout(t, Command().Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	want := "$env:AWS_ACCESS_KEY_ID='fake-access-key-id'\n" +
		"$env:AWS_SECRET_ACCESS_KEY='fake-secret-access-key'\n" +
		"$env:AWS_SESSION_TOKEN='fake-session-token'\n"
	if output != want {
		t.Fatalf("unexpected output:\n got: %q\nwant: %q", output, want)
	}
}

func restoreExportPSCommandState(t *testing.T) {
	t.Helper()

	originalGetSSOCredentials := getSSOCredentials
	t.Cleanup(func() {
		getSSOCredentials = originalGetSSOCredentials
		viper.Reset()
	})

	viper.Set("profile", "dev")
	viper.Set("home-directory", "/fixture-home")
	getSSOCredentials = func(string, string) (*sso.GetRoleCredentialsOutput, *cfg.SSOConfig, *aws.Config, error) {
		creds, ssoConfig, awsConfig := testcreds.FakeCredentialResult()
		return creds, ssoConfig, awsConfig, nil
	}
}
