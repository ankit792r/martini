package flutter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdatePackageNameReplacesContentAndRenamesDirs(t *testing.T) {
	dir := t.TempDir()

	kotlinDir := filepath.Join(dir, "android/app/src/main/kotlin/com/example/oldapp")
	if err := os.MkdirAll(kotlinDir, 0o755); err != nil {
		t.Fatal(err)
	}

	manifestPath := filepath.Join(dir, "android/app/src/main/AndroidManifest.xml")
	if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		manifestPath: `<?xml version="1.0"?>
<manifest package="com.example.oldapp">
    <application android:name="com.example.oldapp.MainApplication" />
</manifest>`,
		filepath.Join(dir, "android/app/build.gradle"): `
android {
    namespace "com.example.oldapp"
    defaultConfig {
        applicationId "com.example.oldapp"
    }
}`,
		filepath.Join(kotlinDir, "MainActivity.kt"): "package com.example.oldapp\n",
		filepath.Join(dir, "ios/Runner/Info.plist"): `<key>CFBundleIdentifier</key>
<string>com.example.oldapp</string>`,
	}

	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result, err := updatePackageName(dir, "com.example.oldapp", "com.example.newapp")
	if err != nil {
		t.Fatal(err)
	}

	if result.FilesUpdated != 4 {
		t.Fatalf("expected 4 updated files, got %d", result.FilesUpdated)
	}
	if result.DirsRenamed != 1 {
		t.Fatalf("expected 1 renamed directory, got %d", result.DirsRenamed)
	}

	newKotlinDir := filepath.Join(dir, "android/app/src/main/kotlin/com/example/newapp")
	if _, err := os.Stat(newKotlinDir); err != nil {
		t.Fatalf("expected renamed kotlin dir: %v", err)
	}

	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(manifest), "com.example.oldapp") {
		t.Fatalf("manifest still contains old package name:\n%s", manifest)
	}
	if !strings.Contains(string(manifest), "com.example.newapp") {
		t.Fatalf("manifest missing new package name:\n%s", manifest)
	}

	mainActivity, err := os.ReadFile(filepath.Join(newKotlinDir, "MainActivity.kt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(mainActivity) != "package com.example.newapp\n" {
		t.Fatalf("unexpected MainActivity content: %q", mainActivity)
	}
}

func TestUpdatePackageNameMissingProjectPath(t *testing.T) {
	_, err := updatePackageName("/path/that/does/not/exist", "com.old.app", "com.new.app")
	if err == nil {
		t.Fatal("expected error for missing project path")
	}
}

func TestUpdatePackageNameRenamesJavaSourceDir(t *testing.T) {
	dir := t.TempDir()

	javaDir := filepath.Join(dir, "android/app/src/main/java/com/example/oldapp")
	if err := os.MkdirAll(javaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(javaDir, "MainActivity.java"), []byte("package com.example.oldapp;\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := updatePackageName(dir, "com.example.oldapp", "com.example.newapp")
	if err != nil {
		t.Fatal(err)
	}
	if result.DirsRenamed != 1 {
		t.Fatalf("expected 1 renamed directory, got %d", result.DirsRenamed)
	}

	newDir := filepath.Join(dir, "android/app/src/main/java/com/example/newapp")
	if _, err := os.Stat(filepath.Join(newDir, "MainActivity.java")); err != nil {
		t.Fatalf("expected java source in renamed dir: %v", err)
	}
}

func TestUpdatePackageNameRenamesWhenParentSegmentChanges(t *testing.T) {
	dir := t.TempDir()

	kotlinDir := filepath.Join(dir, "android/app/src/main/kotlin/com/example/textpert")
	if err := os.MkdirAll(kotlinDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(kotlinDir, "MainActivity.kt"), []byte("package com.example.textpert\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := updatePackageName(dir, "com.example.textpert", "com.system74.textpert")
	if err != nil {
		t.Fatal(err)
	}
	if result.DirsRenamed != 1 {
		t.Fatalf("expected 1 renamed directory, got %d", result.DirsRenamed)
	}

	newDir := filepath.Join(dir, "android/app/src/main/kotlin/com/system74/textpert")
	if _, err := os.Stat(newDir); err != nil {
		t.Fatalf("expected renamed package dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "android/app/src/main/kotlin/com/example")); err == nil {
		t.Fatal("expected old parent dir to be removed")
	} else if !os.IsNotExist(err) {
		t.Fatalf("unexpected stat error: %v", err)
	}
}

func TestUpdatePackageNameSkipsHiddenPaths(t *testing.T) {
	dir := t.TempDir()

	manifest := filepath.Join(dir, "AndroidManifest.xml")
	if err := os.WriteFile(manifest, []byte(`package="com.example.oldapp"`), 0o644); err != nil {
		t.Fatal(err)
	}

	hiddenDir := filepath.Join(dir, ".devbox", "nix", "profile", "default")
	if err := os.MkdirAll(hiddenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hiddenDir, "config"), []byte("com.example.oldapp"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := updatePackageName(dir, "com.example.oldapp", "com.example.newapp")
	if err != nil {
		t.Fatal(err)
	}
	if result.FilesUpdated != 1 {
		t.Fatalf("expected only visible file to update, got %d updates", result.FilesUpdated)
	}

	hiddenConfig, err := os.ReadFile(filepath.Join(hiddenDir, "config"))
	if err != nil {
		t.Fatal(err)
	}
	if string(hiddenConfig) != "com.example.oldapp" {
		t.Fatalf("hidden file should not be modified: %s", hiddenConfig)
	}
}

func TestApplyReplacements(t *testing.T) {
	got := applyReplacements(
		`id "com.example.oldapp" path "com/example/oldapp"`,
		"com.example.oldapp",
		"com.example.newapp",
		"com/example/oldapp",
		"com/example/newapp",
	)
	want := `id "com.example.newapp" path "com/example/newapp"`
	if got != want {
		t.Fatalf("unexpected replacement:\n got: %s\nwant: %s", got, want)
	}
}
