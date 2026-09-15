package flutter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	keystore "github.com/pavlo-v-chernykh/keystore-go/v4"
)

func TestGenerateUploadKeystore(t *testing.T) {
	dir := t.TempDir()
	androidDir := filepath.Join(dir, "android")
	if err := os.MkdirAll(androidDir, 0o755); err != nil {
		t.Fatal(err)
	}

	home := t.TempDir()

	cfg := KeystoreConfig{
		Name:          "test-upload-keystore",
		StorePassword: "storepass123",
		KeyPassword:   "storepass123",
		Alias:         "upload",
		CommonName:    "Test Upload",
	}

	// generateUploadKeystore uses os.UserHomeDir(), not HOME on all platforms consistently,
	// so test createJKS + file writes via lower-level helpers instead.
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

	ks := keystore.New()
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

	credentialsPath := filepath.Join(home, cfg.Name+"-upload-key.txt")
	if err := os.WriteFile(credentialsPath, []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}

	result := &KeystoreResult{
		HomePath:      homePath,
		ProjectPath:   projectPath,
		CredentialsPath: credentialsPath,
		Alias:         cfg.Alias,
		StorePassword: cfg.StorePassword,
		KeyPassword:   cfg.KeyPassword,
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

	if err := applySigningConfig(dir, SigningConfig{
		StoreFile:     result.HomePath,
		StorePassword: result.StorePassword,
		KeyPassword:   result.KeyPassword,
		KeyAlias:      result.Alias,
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
	if !strings.Contains(string(keyProps), "keyAlias=upload") {
		t.Fatalf("unexpected key.properties: %s", keyProps)
	}

	updated, err := os.ReadFile(ktsPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(updated)
	for _, want := range []string{
		"keystoreProperties",
		`signingConfigs {`,
		`signingConfig = signingConfigs.getByName("release")`,
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("missing %q in build.gradle.kts:\n%s", want, content)
		}
	}
}

func TestKeystoreConfigFromAnswers(t *testing.T) {
	cfg := keystoreConfigFromAnswers([]string{"myapp", "secret", "", "upload", "My App"})
	if cfg.Name != "myapp" || cfg.StorePassword != "secret" || cfg.KeyPassword != "secret" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}
