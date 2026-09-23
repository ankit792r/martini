package keystore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesKeystoreAndKeyPropertiesWithoutGradleChange(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := os.WriteFile(filepath.Join(project, "settings.gradle.kts"), []byte("rootProject.name = \"app\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ktsPath := filepath.Join(project, "app", "build.gradle.kts")
	if err := os.MkdirAll(filepath.Dir(ktsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	original := `android {
    namespace = "com.example.app"
}
`
	if err := os.WriteFile(ktsPath, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Run(project, Options{
		KeystoreName:  "upload-keystore",
		StorePassword: "storepass123",
		KeyPassword:   "storepass123",
		KeyAlias:      "upload",
		CommonName:    "Example App",
	})
	if err != nil {
		t.Fatal(err)
	}

	keyProps, err := os.ReadFile(filepath.Join(project, "key.properties"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(keyProps), "storeFile=upload-keystore.jks") {
		t.Fatalf("unexpected key.properties: %s", keyProps)
	}

	updated, err := os.ReadFile(ktsPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(updated) != original {
		t.Fatalf("build.gradle.kts should not change: %s", updated)
	}
}
