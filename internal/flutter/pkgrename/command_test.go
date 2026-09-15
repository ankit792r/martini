package pkgrename

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunMissingAnswers(t *testing.T) {
	if err := Run(t.TempDir(), []string{"only-one"}); err == nil {
		t.Fatal("expected error for missing answers")
	}
}

func TestRunUpdatesManifest(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "AndroidManifest.xml")
	if err := os.WriteFile(manifest, []byte(`package="com.old.app"`), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run(dir, []string{"com.old.app", "com.new.app"}); err != nil {
		t.Fatal(err)
	}

	updated, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if string(updated) != `package="com.new.app"` {
		t.Fatalf("unexpected content: %s", updated)
	}
}
