package flutter

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
		path = expandProjectPath(path)
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("project path: %w", err)
	}

	pubspec := filepath.Join(abs, "pubspec.yaml")
	if _, err := os.Stat(pubspec); err != nil {
		return "", fmt.Errorf("pubspec.yaml not found in %s", abs)
	}

	return abs, nil
}

func expandProjectPath(path string) string {
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
