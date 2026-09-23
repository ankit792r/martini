package flutter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveProjectPathUsesCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte("name: test\n"), 0o644); err != nil {
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

func TestResolveProjectPathExpandsHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	project := filepath.Join(home, "myapp")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "pubspec.yaml"), []byte("name: test\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveProjectPath("~/myapp")
	if err != nil {
		t.Fatal(err)
	}
	if got != project {
		t.Fatalf("expected %s, got %s", project, got)
	}
}

func TestResolveProjectPathMissingPubspec(t *testing.T) {
	dir := t.TempDir()
	if _, err := ResolveProjectPath(dir); err == nil {
		t.Fatal("expected error for missing pubspec.yaml")
	}
}
