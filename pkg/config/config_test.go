package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaxxstorm/aws-sso-creds/internal/testutil"
)

func TestGetSSOConfigDefaultProfile(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteAWSConfig(t, home, `
[default]
sso_start_url = https://example.awsapps.com/start
sso_region = us-west-2
sso_account_id = 123456789012
sso_role_name = AdministratorAccess
`)

	got, err := GetSSOConfig("", home)
	if err != nil {
		t.Fatalf("GetSSOConfig returned error: %v", err)
	}

	assertSSOConfig(t, got, SSOConfig{
		StartURL:  "https://example.awsapps.com/start",
		Region:    "us-west-2",
		AccountID: "123456789012",
		RoleName:  "AdministratorAccess",
	})
}

func TestGetSSOConfigNamedProfile(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteAWSConfig(t, home, `
[profile dev]
sso_start_url = https://dev.awsapps.com/start
sso_region = us-east-1
sso_account_id = 210987654321
sso_role_name = ReadOnly
`)

	got, err := GetSSOConfig("dev", home)
	if err != nil {
		t.Fatalf("GetSSOConfig returned error: %v", err)
	}

	assertSSOConfig(t, got, SSOConfig{
		StartURL:  "https://dev.awsapps.com/start",
		Region:    "us-east-1",
		AccountID: "210987654321",
		RoleName:  "ReadOnly",
	})
}

func TestGetSSOConfigMergesSSOSession(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteAWSConfig(t, home, `
[sso-session corp]
sso_start_url = https://corp.awsapps.com/start
sso_region = eu-west-1

[profile prod]
sso_session = corp
sso_account_id = 111122223333
sso_role_name = PowerUser
`)

	got, err := GetSSOConfig("prod", home)
	if err != nil {
		t.Fatalf("GetSSOConfig returned error: %v", err)
	}

	assertSSOConfig(t, got, SSOConfig{
		StartURL:  "https://corp.awsapps.com/start",
		Region:    "eu-west-1",
		AccountID: "111122223333",
		RoleName:  "PowerUser",
	})
}

func TestGetSSOConfigProfileOverridesSSOSession(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteAWSConfig(t, home, `
[sso-session corp]
sso_start_url = https://corp.awsapps.com/start
sso_region = eu-west-1
sso_account_id = 999999999999
sso_role_name = SessionRole

[profile prod]
sso_session = corp
sso_start_url = https://override.awsapps.com/start
sso_region = ap-southeast-2
sso_account_id = 111122223333
sso_role_name = ProfileRole
`)

	got, err := GetSSOConfig("prod", home)
	if err != nil {
		t.Fatalf("GetSSOConfig returned error: %v", err)
	}

	assertSSOConfig(t, got, SSOConfig{
		StartURL:  "https://override.awsapps.com/start",
		Region:    "ap-southeast-2",
		AccountID: "111122223333",
		RoleName:  "ProfileRole",
	})
}

func TestGetSSOConfigMissingProfile(t *testing.T) {
	home := testutil.AWSHome(t)
	testutil.WriteAWSConfig(t, home, `[profile dev]`)

	_, err := GetSSOConfig("missing", home)
	if err == nil || err.Error() != "unable to find profile: missing" {
		t.Fatalf("expected missing profile error, got %v", err)
	}
}

func TestGetSSOConfigMissingRequiredSettings(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		profile     string
		wantErrPart string
	}{
		{
			name: "missing session region",
			config: `
[sso-session corp]
sso_start_url = https://corp.awsapps.com/start

[profile prod]
sso_session = corp
sso_account_id = 111122223333
sso_role_name = Admin
`,
			profile:     "prod",
			wantErrPart: `no sso_region in sso-session "corp"`,
		},
		{
			name: "missing profile start URL",
			config: `
[profile dev]
sso_region = us-west-2
sso_account_id = 123456789012
sso_role_name = Admin
`,
			profile:     "dev",
			wantErrPart: `no sso_start_url in profile "dev" and its sso_session`,
		},
		{
			name: "missing account ID",
			config: `
[profile dev]
sso_start_url = https://dev.awsapps.com/start
sso_region = us-west-2
sso_role_name = Admin
`,
			profile:     "dev",
			wantErrPart: `no sso_account_id in profile "dev" or its sso_session`,
		},
		{
			name: "missing role name",
			config: `
[profile dev]
sso_start_url = https://dev.awsapps.com/start
sso_region = us-west-2
sso_account_id = 123456789012
`,
			profile:     "dev",
			wantErrPart: `no sso_role_name in profile "dev" or its sso_session`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := testutil.AWSHome(t)
			testutil.WriteAWSConfig(t, home, tt.config)

			_, err := GetSSOConfig(tt.profile, home)
			if err == nil || !strings.Contains(err.Error(), tt.wantErrPart) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErrPart, err)
			}
		})
	}
}

func assertSSOConfig(t *testing.T, got *SSOConfig, want SSOConfig) {
	t.Helper()

	if *got != want {
		t.Fatalf("unexpected config:\n got: %#v\nwant: %#v", *got, want)
	}
}

func readCacheEntries(t *testing.T, home string) []os.DirEntry {
	t.Helper()

	files, err := os.ReadDir(filepath.Join(home, ".aws", "sso", "cache"))
	if err != nil {
		t.Fatalf("read cache entries: %v", err)
	}
	return files
}
