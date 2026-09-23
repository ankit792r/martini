package keystore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"martini/internal/flutter/gradle"
	shared "martini/internal/keystore"
)

func TestFlutterKeystoreUpdatesGradle(t *testing.T) {
	dir := t.TempDir()
	androidDir := filepath.Join(dir, "android")
	if err := os.MkdirAll(filepath.Join(androidDir, "app"), 0o755); err != nil {
		t.Fatal(err)
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := shared.Config{
		Name:          "test-upload-keystore",
		StorePassword: "storepass123",
		KeyPassword:   "storepass123",
		Alias:         "upload",
		CommonName:    "Test Upload",
	}

	result, err := shared.Generate(androidDir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	ktsPath := filepath.Join(androidDir, "app", "build.gradle.kts")
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
		StoreFile:     result.HomePath,
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
