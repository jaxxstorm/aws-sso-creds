package get

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/internal/testcreds"
	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/spf13/viper"
)

func TestCommandJSONOutput(t *testing.T) {
	restoreGetCommandState(t)

	command := Command()
	command.SetArgs([]string{"--json"})

	output, err := testutil.CaptureStdout(t, command.Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	var got JSON
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("unmarshal output %q: %v", output, err)
	}
	if got.AwsAccessKeyID != testcreds.FakeAccessKeyID {
		t.Fatalf("unexpected access key %q", got.AwsAccessKeyID)
	}
	if got.AwsSecretAccessKey != testcreds.FakeSecretAccessKey {
		t.Fatalf("unexpected secret key %q", got.AwsSecretAccessKey)
	}
	if got.SessionToken != testcreds.FakeSessionToken {
		t.Fatalf("unexpected session token %q", got.SessionToken)
	}
	if !got.ExpireAt.Equal(time.UnixMilli(testcreds.FakeExpiration)) {
		t.Fatalf("unexpected expiration %s", got.ExpireAt)
	}
}

func TestCommandHumanOutput(t *testing.T) {
	restoreGetCommandState(t)

	output, err := testutil.CaptureStdout(t, Command().Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	for _, want := range []string{
		"Your temporary credentials for account",
		testcreds.FakeAccountID,
		"AWS_ACCESS_KEY_ID",
		testcreds.FakeAccessKeyID,
		"AWS_SECRET_ACCESS_KEY",
		testcreds.FakeSecretAccessKey,
		"AWS_SESSION_TOKEN",
		testcreds.FakeSessionToken,
		"These credentials will expire at:",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, output)
		}
	}
}

func restoreGetCommandState(t *testing.T) {
	t.Helper()

	originalGetSSOCredentials := getSSOCredentials
	t.Cleanup(func() {
		getSSOCredentials = originalGetSSOCredentials
		viper.Reset()
	})

	viper.Set("profile", "dev")
	viper.Set("home-directory", "/fixture-home")
	getSSOCredentials = func(profile string, homedir string) (*sso.GetRoleCredentialsOutput, *cfg.SSOConfig, *aws.Config, error) {
		if profile != "dev" {
			t.Fatalf("unexpected profile %q", profile)
		}
		if homedir != "/fixture-home" {
			t.Fatalf("unexpected home directory %q", homedir)
		}
		creds, ssoConfig, awsConfig := testcreds.FakeCredentialResult()
		return creds, ssoConfig, awsConfig, nil
	}
}
