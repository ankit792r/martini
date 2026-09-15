package signing

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"martini/internal/flutter/gradle"
)

func Resolve(projectPath string) (gradle.Properties, error) {
	keyPropsPath := filepath.Join(projectPath, "android", "key.properties")
	if cfg, err := readKeyProperties(keyPropsPath); err == nil {
		return cfg, nil
	}

	for _, credPath := range candidateCredentialFiles(projectPath) {
		if cfg, err := configFromCredentialsFile(credPath); err == nil {
			return cfg, nil
		}
	}

	return gradle.Properties{}, fmt.Errorf("signing config not found: expected android/key.properties or a credentials file in home (e.g. ~/upload-keystore-upload-key.txt)")
}

func Apply(projectPath string, cfg gradle.Properties) error {
	return gradle.Apply(projectPath, cfg)
}

func candidateCredentialFiles(projectPath string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	seen := map[string]bool{}
	var paths []string

	add := func(path string) {
		path = strings.TrimSpace(path)
		if path == "" || seen[path] {
			return
		}
		seen[path] = true
		paths = append(paths, path)
	}

	add(filepath.Join(home, "upload-keystore-upload-key.txt"))

	androidDir := filepath.Join(projectPath, "android")
	matches, _ := filepath.Glob(filepath.Join(androidDir, "*.jks"))
	for _, keystorePath := range matches {
		base := strings.TrimSuffix(filepath.Base(keystorePath), ".jks")
		add(filepath.Join(home, base+"-upload-key.txt"))
	}

	return paths
}

func readKeyProperties(path string) (gradle.Properties, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return gradle.Properties{}, err
	}

	cfg := gradle.Properties{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "storePassword":
			cfg.StorePassword = strings.TrimSpace(value)
		case "keyPassword":
			cfg.KeyPassword = strings.TrimSpace(value)
		case "keyAlias":
			cfg.KeyAlias = strings.TrimSpace(value)
		case "storeFile":
			cfg.StoreFile = strings.TrimSpace(strings.ReplaceAll(value, "\\\\", "\\"))
		}
	}

	if err := validateConfig(cfg); err != nil {
		return gradle.Properties{}, err
	}
	return cfg, nil
}

func configFromCredentialsFile(path string) (gradle.Properties, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return gradle.Properties{}, fmt.Errorf("empty credentials path")
	}

	file, err := os.Open(path)
	if err != nil {
		return gradle.Properties{}, err
	}
	defer file.Close()

	cfg := gradle.Properties{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "Keystore (home):"):
			cfg.StoreFile = strings.TrimSpace(strings.TrimPrefix(line, "Keystore (home):"))
		case strings.HasPrefix(line, "Store password:"):
			cfg.StorePassword = strings.TrimSpace(strings.TrimPrefix(line, "Store password:"))
		case strings.HasPrefix(line, "Key password:"):
			cfg.KeyPassword = strings.TrimSpace(strings.TrimPrefix(line, "Key password:"))
		case strings.HasPrefix(line, "Key alias:"):
			cfg.KeyAlias = strings.TrimSpace(strings.TrimPrefix(line, "Key alias:"))
		}
	}
	if err := scanner.Err(); err != nil {
		return gradle.Properties{}, err
	}

	if cfg.StoreFile == "" {
		return gradle.Properties{}, fmt.Errorf("credentials file missing keystore path")
	}
	if cfg.KeyPassword == "" {
		cfg.KeyPassword = cfg.StorePassword
	}
	if cfg.KeyAlias == "" {
		cfg.KeyAlias = gradle.DefaultUploadAlias
	}

	if err := validateConfig(cfg); err != nil {
		return gradle.Properties{}, err
	}
	return cfg, nil
}

func validateConfig(cfg gradle.Properties) error {
	if strings.TrimSpace(cfg.StoreFile) == "" {
		return fmt.Errorf("store file is required")
	}
	if strings.TrimSpace(cfg.StorePassword) == "" {
		return fmt.Errorf("store password is required")
	}
	if strings.TrimSpace(cfg.KeyPassword) == "" {
		return fmt.Errorf("key password is required")
	}
	if strings.TrimSpace(cfg.KeyAlias) == "" {
		return fmt.Errorf("key alias is required")
	}
	if _, err := os.Stat(cfg.StoreFile); err != nil {
		return fmt.Errorf("keystore file not found: %s", cfg.StoreFile)
	}
	return nil
}
