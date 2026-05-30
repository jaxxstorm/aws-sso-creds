package version

import (
	"errors"
	"strings"
	"testing"

	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
	appversion "github.com/jaxxstorm/aws-sso-creds/pkg/version"
)

func TestCommandUsesLinkedVersion(t *testing.T) {
	restoreVersionCommandState(t)

	appversion.Version = "1.2.3"
	calculateFallbackVersion = func() (string, error) {
		t.Fatal("fallback version calculation should not be called")
		return "", nil
	}

	output, err := testutil.CaptureStdout(t, Command().Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if output != "1.2.3\n" {
		t.Fatalf("unexpected output %q", output)
	}
}

func TestCommandUsesFallbackVersion(t *testing.T) {
	restoreVersionCommandState(t)

	calculateFallbackVersion = func() (string, error) {
		return "0.1.0-alpha.123+abcdef12", nil
	}

	output, err := testutil.CaptureStdout(t, Command().Execute)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if output != "0.1.0-alpha.123+abcdef12\n" {
		t.Fatalf("unexpected output %q", output)
	}
}

func TestCommandReturnsFallbackVersionError(t *testing.T) {
	restoreVersionCommandState(t)

	calculateFallbackVersion = func() (string, error) {
		return "", errors.New("no repository")
	}

	output, err := testutil.CaptureStdout(t, Command().Execute)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "no repository") {
		t.Fatalf("unexpected error %q", err)
	}
	if output != "" {
		t.Fatalf("unexpected output %q", output)
	}
}

func TestCommandHelpNamesAwsSSOCreds(t *testing.T) {
	command := Command()

	if !strings.Contains(command.Long, "aws-sso-creds") {
		t.Fatalf("help text should mention aws-sso-creds: %q", command.Long)
	}
	if strings.Contains(command.Long, "pulumictl") {
		t.Fatalf("help text should not mention pulumictl: %q", command.Long)
	}
}

func restoreVersionCommandState(t *testing.T) {
	t.Helper()

	originalVersion := appversion.Version
	originalCalculateFallbackVersion := calculateFallbackVersion
	appversion.Version = ""

	t.Cleanup(func() {
		appversion.Version = originalVersion
		calculateFallbackVersion = originalCalculateFallbackVersion
	})
}
