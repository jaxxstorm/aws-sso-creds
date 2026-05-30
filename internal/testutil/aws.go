package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func AWSHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	MkdirAll(t, filepath.Join(home, ".aws", "sso", "cache"))
	return home
}

func MkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("create directory %s: %v", path, err)
	}
}

func WriteAWSConfig(t *testing.T, home string, content string) {
	t.Helper()
	WriteFile(t, filepath.Join(home, ".aws", "config"), content)
}

func WriteAWSCredentials(t *testing.T, home string, content string) {
	t.Helper()
	WriteFile(t, filepath.Join(home, ".aws", "credentials"), content)
}

func WriteSSOCache(t *testing.T, home string, name string, content string) {
	t.Helper()
	WriteFile(t, filepath.Join(home, ".aws", "sso", "cache", name), content)
}

func WriteFile(t *testing.T, path string, content string) {
	t.Helper()

	MkdirAll(t, filepath.Dir(path))
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}
