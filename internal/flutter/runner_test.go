package flutter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandsFromFlags(t *testing.T) {
	flags := map[string]bool{
		"update-package-name":      true,
		"generate-upload-keystore": false,
		"update-signing-config":    true,
	}

	got := CommandsFromFlags(flags)
	if len(got) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(got))
	}
	if got[0] != CmdUpdatePackageName || got[1] != CmdUpdateSigningConfig {
		t.Fatalf("unexpected command order: %v", got)
	}
}

func TestRunUpdatePackageName(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "AndroidManifest.xml")
	content := `<manifest package="com.old.app" />`
	if err := os.WriteFile(manifest, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	session := NewSession(dir)
	session.SetAnswers(CmdUpdatePackageName, []string{"com.old.app", "com.new.app"})

	if err := runUpdatePackageName(session); err != nil {
		t.Fatal(err)
	}

	updated, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if string(updated) != `<manifest package="com.new.app" />` {
		t.Fatalf("unexpected manifest content: %s", updated)
	}
}

func TestRunUpdatePackageNameMissingAnswers(t *testing.T) {
	session := NewSession(t.TempDir())
	if err := runUpdatePackageName(session); err == nil {
		t.Fatal("expected error for missing answers")
	}
}
