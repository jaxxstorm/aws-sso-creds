package set

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/internal/testcreds"
	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/spf13/viper"
)

func TestCommandWritesTemporaryCredentialsToFixtureFiles(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteAWSConfig(t, home, `
[profile dev]
region = us-west-2
`)
	testutil.WriteAWSCredentials(t, home, `
[default]
aws_access_key_id = existing
`)
	restoreSetCommandState(t, home)

	command := Command()
	command.SetArgs([]string{"generated"})
	output, err := testutil.CaptureStdout(t, command.Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	credsPath := filepath.Join(home, ".aws", "credentials")
	configPath := filepath.Join(home, ".aws", "config")
	if !strings.Contains(output, credsPath) || !strings.Contains(output, configPath) {
		t.Fatalf("expected output to contain fixture paths, got:\n%s", output)
	}
	if !strings.Contains(output, "credentials saved to profile: generated") {
		t.Fatalf("expected save message, got:\n%s", output)
	}

	credentialsContent := readFile(t, credsPath)
	for _, want := range []string{
		"[generated]",
		"aws_access_key_id = " + testcreds.FakeAccessKeyID,
		"aws_secret_access_key = " + testcreds.FakeSecretAccessKey,
		"aws_session_token = " + testcreds.FakeSessionToken,
	} {
		if !strings.Contains(credentialsContent, want) {
			t.Fatalf("expected credentials file to contain %q, got:\n%s", want, credentialsContent)
		}
	}

	configContent := readFile(t, configPath)
	if !strings.Contains(configContent, "[profile generated]") {
		t.Fatalf("expected config file to contain generated profile, got:\n%s", configContent)
	}
}

func restoreSetCommandState(t *testing.T, home string) {
	t.Helper()

	originalGetSSOCredentials := getSSOCredentials
	t.Cleanup(func() {
		getSSOCredentials = originalGetSSOCredentials
		viper.Reset()
	})

	viper.Set("profile", "dev")
	viper.Set("home-directory", home)
	getSSOCredentials = func(profile string, homedir string) (*sso.GetRoleCredentialsOutput, *cfg.SSOConfig, *aws.Config, error) {
		if profile != "dev" {
			t.Fatalf("unexpected profile %q", profile)
		}
		if homedir != home {
			t.Fatalf("unexpected home directory %q", homedir)
		}
		creds, ssoConfig, awsConfig := testcreds.FakeCredentialResult()
		return creds, ssoConfig, awsConfig, nil
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	return string(content)
}
