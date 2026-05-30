package exec

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/internal/testcreds"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/spf13/viper"
)

func TestCommandBuildsExpectedEnvironment(t *testing.T) {
	restoreExecCommandState(t)

	var gotBinary string
	var gotArgs []string
	var gotEnv []string
	lookPath = func(file string) (string, error) {
		if file != "aws" {
			t.Fatalf("unexpected path lookup %q", file)
		}
		return "/usr/bin/aws", nil
	}
	execCommand = func(binary string, args []string, env []string) error {
		gotBinary = binary
		gotArgs = append([]string(nil), args...)
		gotEnv = append([]string(nil), env...)
		return nil
	}

	command := Command()
	command.SetArgs([]string{"aws", "s3", "ls"})
	if err := command.Execute(); err != nil {
		t.Fatalf("execute command: %v", err)
	}

	if gotBinary != "/usr/bin/aws" {
		t.Fatalf("unexpected binary %q", gotBinary)
	}
	if len(gotArgs) != 3 || gotArgs[0] != "aws" || gotArgs[1] != "s3" || gotArgs[2] != "ls" {
		t.Fatalf("unexpected args %#v", gotArgs)
	}
	for _, want := range []string{
		"AWS_ACCESS_KEY_ID=" + testcreds.FakeAccessKeyID,
		"AWS_SECRET_ACCESS_KEY=" + testcreds.FakeSecretAccessKey,
		"AWS_SESSION_TOKEN=" + testcreds.FakeSessionToken,
		"AWS_DEFAULT_REGION=" + testcreds.FakeRegion,
	} {
		if !envContains(gotEnv, want) {
			t.Fatalf("expected env to contain %q, got %#v", want, gotEnv)
		}
	}
}

func restoreExecCommandState(t *testing.T) {
	t.Helper()

	originalGetSSOCredentials := getSSOCredentials
	originalLookPath := lookPath
	originalExecCommand := execCommand
	originalEnv := os.Environ()
	t.Cleanup(func() {
		getSSOCredentials = originalGetSSOCredentials
		lookPath = originalLookPath
		execCommand = originalExecCommand
		viper.Reset()
		os.Clearenv()
		for _, entry := range originalEnv {
			key, value, ok := strings.Cut(entry, "=")
			if ok {
				_ = os.Setenv(key, value)
			}
		}
	})

	os.Clearenv()
	viper.Set("profile", "dev")
	viper.Set("home-directory", "/fixture-home")
	getSSOCredentials = func(string, string) (*sso.GetRoleCredentialsOutput, *cfg.SSOConfig, *aws.Config, error) {
		creds, ssoConfig, awsConfig := testcreds.FakeCredentialResult()
		return creds, ssoConfig, awsConfig, nil
	}
}

func envContains(env []string, want string) bool {
	for _, entry := range env {
		if entry == want {
			return true
		}
	}
	return false
}
