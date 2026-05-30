package helper

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/internal/testcreds"
	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/spf13/viper"
)

func TestCommandOutput(t *testing.T) {
	restoreHelperCommandState(t)

	output, err := testutil.CaptureStdout(t, Command().Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	var got CredentialsProcessOutput
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("unmarshal output %q: %v", output, err)
	}
	if got.Version != 1 {
		t.Fatalf("unexpected version %d", got.Version)
	}
	if got.AccessKeyID != testcreds.FakeAccessKeyID {
		t.Fatalf("unexpected access key %q", got.AccessKeyID)
	}
	if got.SecretAccessKey != testcreds.FakeSecretAccessKey {
		t.Fatalf("unexpected secret key %q", got.SecretAccessKey)
	}
	if got.SessionToken != testcreds.FakeSessionToken {
		t.Fatalf("unexpected session token %q", got.SessionToken)
	}
	if got.Expiration != time.Unix(testcreds.FakeExpiration/1000, 0).Format(time.RFC3339) {
		t.Fatalf("unexpected expiration %q", got.Expiration)
	}
}

func restoreHelperCommandState(t *testing.T) {
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
