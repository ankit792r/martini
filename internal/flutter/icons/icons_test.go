package icons

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateCopiesIconKitchenOutput(t *testing.T) {
	iconsDir := t.TempDir()
	projectDir := t.TempDir()

	writeFile(t, filepath.Join(projectDir, "pubspec.yaml"), "name: test\n")
	writeFile(t, filepath.Join(projectDir, "android", "app", "src", "main", "res", "mipmap-mdpi", "ic_launcher.png"), "old-android")
	writeFile(t, filepath.Join(projectDir, "ios", "Runner", "Assets.xcassets", "AppIcon.appiconset", "Contents.json"), `{}`)
	writeFile(t, filepath.Join(projectDir, "ios", "Runner", "Assets.xcassets", "AppIcon.appiconset", "Icon-App-60x60@2x.png"), "old-ios")
	writeFile(t, filepath.Join(projectDir, "web", "manifest.json"), `{}`)
	writeFile(t, filepath.Join(projectDir, "web", "icons", "Icon-192.png"), "old-web")

	writeFile(t, filepath.Join(iconsDir, "android", "res", "mipmap-mdpi", "ic_launcher.png"), "new-android")
	writeFile(t, filepath.Join(iconsDir, "android", "res", "mipmap-anydpi-v26", "ic_launcher.xml"), "<adaptive/>")
	writeFile(t, filepath.Join(iconsDir, "ios", "AppIcon@2x.png"), "new-ios")
	writeFile(t, filepath.Join(iconsDir, "web", "icon-192.png"), "new-web")
	writeFile(t, filepath.Join(iconsDir, "web", "icon-512.png"), "new-web-512")
	writeFile(t, filepath.Join(iconsDir, "web", "favicon.ico"), "new-favicon")

	result, err := Update(projectDir, iconsDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.AndroidFiles != 2 {
		t.Fatalf("expected 2 android files, got %d", result.AndroidFiles)
	}
	if result.IOSFiles != 1 {
		t.Fatalf("expected 1 ios file, got %d", result.IOSFiles)
	}
	if result.WebFiles < 3 {
		t.Fatalf("expected at least 3 web files, got %d", result.WebFiles)
	}

	assertFileContents(t, filepath.Join(projectDir, "android", "app", "src", "main", "res", "mipmap-mdpi", "ic_launcher.png"), "new-android")
	assertFileContents(t, filepath.Join(projectDir, "android", "app", "src", "main", "res", "mipmap-anydpi-v26", "ic_launcher.xml"), "<adaptive/>")
	assertFileContents(t, filepath.Join(projectDir, "ios", "Runner", "Assets.xcassets", "AppIcon.appiconset", "Icon-App-60x60@2x.png"), "new-ios")
	assertFileContents(t, filepath.Join(projectDir, "web", "icons", "Icon-192.png"), "new-web")
	assertFileContents(t, filepath.Join(projectDir, "web", "favicon.png"), "new-web")
}

func TestResolveIconsPathExpandsHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := filepath.Join(home, "IconKitchen-Output")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveIconsPath("~/IconKitchen-Output")
	if err != nil {
		t.Fatal(err)
	}
	if got != dir {
		t.Fatalf("expected %s, got %s", dir, got)
	}
}

func TestUpdateMissingProject(t *testing.T) {
	if _, err := Update(t.TempDir(), t.TempDir()); err == nil {
		t.Fatal("expected error for missing pubspec.yaml")
	}
}

func TestResolveIconsPathRequiresPath(t *testing.T) {
	if _, err := ResolveIconsPath(""); err == nil {
		t.Fatal("expected error for empty icons path")
	}
}

func TestRunRequiresIconsPath(t *testing.T) {
	if err := Run(t.TempDir(), nil); err == nil {
		t.Fatal("expected error for missing icons path")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertFileContents(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s: expected %q, got %q", path, want, got)
	}
}
