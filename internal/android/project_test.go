package android

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveProjectPathUsesCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "settings.gradle.kts"), []byte("rootProject.name = \"test\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	want, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ResolveProjectPath("")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolveProjectPathRejectsNonAndroidProject(t *testing.T) {
	dir := t.TempDir()
	if _, err := ResolveProjectPath(dir); err == nil {
		t.Fatal("expected error for non-android directory")
	}
}
