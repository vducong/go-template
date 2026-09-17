package cfg

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestLoad_DoesNotPrintSecrets(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: test\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AUTH_JWT_SECRET", "jwt-s3cret")
	t.Setenv("AUTH_API_KEYS", "api-key-s3cret")

	var configs *Config
	out := captureStdout(t, func() {
		var err error
		if configs, err = Load(path); err != nil {
			t.Fatal(err)
		}
	})

	if configs.Auth.JWTSecret != "jwt-s3cret" || !slices.Contains(configs.Auth.APIKeys, "api-key-s3cret") {
		t.Fatalf("secrets must be loaded, or the check below is vacuous: %+v", configs.Auth)
	}
	for _, secret := range []string{"jwt-s3cret", "api-key-s3cret"} {
		if strings.Contains(out, secret) {
			t.Errorf("stdout contains %q:\n%s", secret, out)
		}
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}
