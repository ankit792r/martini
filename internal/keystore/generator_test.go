package keystore

import (
	"os"
	"strings"
	"testing"

	jks "github.com/pavlo-v-chernykh/keystore-go/v4"
)

func TestGenerateJKS(t *testing.T) {
	projectStore := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := Config{
		Name:          "test-upload-keystore",
		StorePassword: "storepass123",
		KeyPassword:   "storepass123",
		Alias:         "upload",
		CommonName:    "Test Upload",
	}

	result, err := Generate(projectStore, cfg)
	if err != nil {
		t.Fatal(err)
	}

	ks := jks.New()
	f, err := os.Open(result.HomePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := ks.Load(f, []byte(cfg.StorePassword)); err != nil {
		t.Fatal(err)
	}
	if !ks.IsPrivateKeyEntry(cfg.Alias) {
		t.Fatal("expected private key entry in keystore")
	}

	if _, err := os.Stat(result.ProjectPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(result.CredentialsPath); err != nil {
		t.Fatal(err)
	}
}

func TestConfigFromAnswers(t *testing.T) {
	cfg := ConfigFromAnswers([]string{"myapp", "secret", "", "upload", "My App"})
	if cfg.Name != "myapp" || cfg.StorePassword != "secret" || cfg.KeyPassword != "secret" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestGenerateRejectsExistingKeystore(t *testing.T) {
	projectStore := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := Config{
		Name:          "upload-keystore",
		StorePassword: "pass",
		KeyPassword:   "pass",
		Alias:         "upload",
		CommonName:    "Test",
	}

	if _, err := Generate(projectStore, cfg); err != nil {
		t.Fatal(err)
	}
	_, err := Generate(projectStore, cfg)
	if err == nil {
		t.Fatal("expected error when keystore already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("unexpected error: %v", err)
	}
}
