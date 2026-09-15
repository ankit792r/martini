package flutter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteReleaseOutputsMarkdown(t *testing.T) {
	dir := t.TempDir()
	artifacts := []ReleaseArtifact{
		{Name: "app-arm64-v8a-release.apk", Type: "release", SHA1: "abc123"},
		{Name: "app-release.aab", Type: "release", SHA1: "def456"},
	}

	outputPath := filepath.Join(dir, "release-outputs.md")
	if err := writeReleaseOutputsMarkdown(outputPath, artifacts); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	for _, want := range []string{
		"### Release outputs",
		"| app-arm64-v8a-release.apk | release | `abc123` |",
		"| app-release.aab | release | `def456` |",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in:\n%s", want, text)
		}
	}
}

func TestCollectReleaseArtifacts(t *testing.T) {
	dir := t.TempDir()
	apkDir := filepath.Join(dir, "build", "app", "outputs", "flutter-apk")
	aabDir := filepath.Join(dir, "build", "app", "outputs", "bundle", "release")
	if err := os.MkdirAll(apkDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(aabDir, 0o755); err != nil {
		t.Fatal(err)
	}

	apkPath := filepath.Join(apkDir, "demo-arm64-v8a-release.apk")
	aabPath := filepath.Join(aabDir, "app-release.aab")
	if err := os.WriteFile(apkPath, []byte("apk-content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(aabPath, []byte("aab-content"), 0o644); err != nil {
		t.Fatal(err)
	}

	artifacts, err := collectReleaseArtifacts(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(artifacts))
	}
}
