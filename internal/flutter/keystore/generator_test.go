package keystore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"martini/internal/flutter/gradle"

	jks "github.com/pavlo-v-chernykh/keystore-go/v4"
)

func TestGenerateJKS(t *testing.T) {
	dir := t.TempDir()
	androidDir := filepath.Join(dir, "android")
	if err := os.MkdirAll(androidDir, 0o755); err != nil {
		t.Fatal(err)
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := Config{
		Name:          "test-upload-keystore",
		StorePassword: "storepass123",
		KeyPassword:   "storepass123",
		Alias:         "upload",
		CommonName:    "Test Upload",
	}

	keystoreBytes, _, err := createJKS(cfg)
	if err != nil {
		t.Fatal(err)
	}

	homePath := filepath.Join(home, cfg.Name+".jks")
	projectPath := filepath.Join(androidDir, cfg.Name+".jks")
	if err := os.WriteFile(homePath, keystoreBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(homePath, projectPath); err != nil {
		t.Fatal(err)
	}

	ks := jks.New()
	f, err := os.Open(homePath)
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

	ktsPath := filepath.Join(androidDir, "app", "build.gradle.kts")
	if err := os.MkdirAll(filepath.Dir(ktsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	sample := `plugins {
    id("com.android.application")
}

android {
    namespace = "com.example.test"
    defaultConfig {
        applicationId = "com.example.test"
    }
    buildTypes {
        release {
            signingConfig = signingConfigs.getByName("debug")
        }
    }
}
`
	if err := os.WriteFile(ktsPath, []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := gradle.Apply(dir, gradle.Properties{
		StoreFile:     homePath,
		StorePassword: cfg.StorePassword,
		KeyPassword:   cfg.KeyPassword,
		KeyAlias:      cfg.Alias,
	}); err != nil {
		t.Fatal(err)
	}

	keyProps, err := os.ReadFile(filepath.Join(androidDir, "key.properties"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(keyProps), "storePassword=storepass123") {
		t.Fatalf("unexpected key.properties: %s", keyProps)
	}

	updated, err := os.ReadFile(ktsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(updated), "keystoreProperties") {
		t.Fatalf("missing gradle signing config: %s", updated)
	}
}

func TestConfigFromAnswers(t *testing.T) {
	cfg := configFromAnswers([]string{"myapp", "secret", "", "upload", "My App"})
	if cfg.Name != "myapp" || cfg.StorePassword != "secret" || cfg.KeyPassword != "secret" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}
