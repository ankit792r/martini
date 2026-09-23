package android

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveProjectPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("project path: %w", err)
		}
		path = cwd
	} else {
		path = expandHome(path)
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("project path: %w", err)
	}

	if !isAndroidProject(abs) {
		return "", fmt.Errorf("not an Android Studio project (expected settings.gradle(.kts) or app/build.gradle(.kts) in %s)", abs)
	}

	return abs, nil
}

func isAndroidProject(root string) bool {
	markers := []string{
		"settings.gradle.kts",
		"settings.gradle",
		filepath.Join("app", "build.gradle.kts"),
		filepath.Join("app", "build.gradle"),
	}
	for _, marker := range markers {
		if _, err := os.Stat(filepath.Join(root, marker)); err == nil {
			return true
		}
	}
	return false
}

func expandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		if path == "~" {
			return home
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return path
}
