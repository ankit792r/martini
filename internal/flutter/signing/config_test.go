package signing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveFromKeyProperties(t *testing.T) {
	dir := t.TempDir()
	androidDir := filepath.Join(dir, "android")
	if err := os.MkdirAll(androidDir, 0o755); err != nil {
		t.Fatal(err)
	}

	keystorePath := filepath.Join(dir, "test.jks")
	if err := os.WriteFile(keystorePath, []byte("keystore"), 0o600); err != nil {
		t.Fatal(err)
	}

	keyProps := "storePassword=store123\nkeyPassword=key123\nkeyAlias=upload\nstoreFile=" + keystorePath + "\n"
	if err := os.WriteFile(filepath.Join(androidDir, "key.properties"), []byte(keyProps), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Resolve(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.StorePassword != "store123" || cfg.KeyAlias != "upload" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestResolveFromCredentialsFile(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	androidDir := filepath.Join(dir, "android")
	if err := os.MkdirAll(androidDir, 0o755); err != nil {
		t.Fatal(err)
	}

	keystorePath := filepath.Join(home, "upload-keystore.jks")
	if err := os.WriteFile(keystorePath, []byte("keystore"), 0o600); err != nil {
		t.Fatal(err)
	}

	creds := "" +
		"Keystore (home): " + keystorePath + "\n" +
		"Store password: secret\n" +
		"Key password: secret\n" +
		"Key alias: upload\n"
	if err := os.WriteFile(filepath.Join(home, "upload-keystore-upload-key.txt"), []byte(creds), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Resolve(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.StoreFile != keystorePath {
		t.Fatalf("unexpected store file: %s", cfg.StoreFile)
	}
}

func TestRunUpdatesGradle(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	androidDir := filepath.Join(dir, "android", "app")
	if err := os.MkdirAll(androidDir, 0o755); err != nil {
		t.Fatal(err)
	}

	keystorePath := filepath.Join(home, "upload-keystore.jks")
	if err := os.WriteFile(keystorePath, []byte("keystore"), 0o600); err != nil {
		t.Fatal(err)
	}

	creds := "" +
		"Keystore (home): " + keystorePath + "\n" +
		"Store password: secret\n" +
		"Key password: secret\n" +
		"Key alias: upload\n"
	if err := os.WriteFile(filepath.Join(home, "upload-keystore-upload-key.txt"), []byte(creds), 0o600); err != nil {
		t.Fatal(err)
	}

	kts := `plugins {
    id("com.android.application")
}

android {
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
	if err := os.WriteFile(filepath.Join(androidDir, "build.gradle.kts"), []byte(kts), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run(dir); err != nil {
		t.Fatal(err)
	}

	keyProps, err := os.ReadFile(filepath.Join(dir, "android", "key.properties"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(keyProps), "storePassword=") {
		t.Fatalf("missing key.properties content: %s", keyProps)
	}

	updated, err := os.ReadFile(filepath.Join(androidDir, "build.gradle.kts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(updated), "keystoreProperties") {
		t.Fatalf("gradle not updated: %s", updated)
	}
}
