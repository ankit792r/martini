package flutter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SigningConfig struct {
	StorePassword string
	KeyPassword   string
	KeyAlias      string
	StoreFile     string
}

func runUpdateSigningConfig(session *Session) error {
	fmt.Printf("Updating signing config for %s\n", session.ProjectPath)

	cfg, err := resolveSigningConfig(session.ProjectPath)
	if err != nil {
		return err
	}

	if err := applySigningConfig(session.ProjectPath, cfg); err != nil {
		return err
	}

	fmt.Printf("  key.properties: %s/android/key.properties\n", session.ProjectPath)
	fmt.Printf("  store file: %s\n", cfg.StoreFile)
	fmt.Println("  build.gradle updated for release signing")
	fmt.Println("  status: signing config updated")

	return nil
}

func resolveSigningConfig(projectPath string) (SigningConfig, error) {
	keyPropsPath := filepath.Join(projectPath, "android", "key.properties")
	if cfg, err := readKeyProperties(keyPropsPath); err == nil {
		return cfg, nil
	}

	for _, credPath := range candidateCredentialFiles(projectPath) {
		if cfg, err := signingConfigFromCredentialsFile(credPath); err == nil {
			return cfg, nil
		}
	}

	return SigningConfig{}, fmt.Errorf("signing config not found: expected android/key.properties or a credentials file in home (e.g. ~/upload-keystore-upload-key.txt)")
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

func applySigningConfig(projectPath string, cfg SigningConfig) error {
	keyPropertiesPath := filepath.Join(projectPath, "android", "key.properties")
	result := &KeystoreResult{
		HomePath:      cfg.StoreFile,
		Alias:         cfg.KeyAlias,
		StorePassword: cfg.StorePassword,
		KeyPassword:   cfg.KeyPassword,
	}
	if err := writeKeyProperties(keyPropertiesPath, result); err != nil {
		return err
	}
	return updateGradleSigning(projectPath)
}

func updateGradleSigning(projectPath string) error {
	androidDir := filepath.Join(projectPath, "android")
	ktsPath := filepath.Join(androidDir, "app", "build.gradle.kts")
	if _, err := os.Stat(ktsPath); err == nil {
		return updateBuildGradleKTS(ktsPath)
	}

	gradlePath := filepath.Join(androidDir, "app", "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		return updateBuildGradleGroovy(gradlePath)
	}

	return fmt.Errorf("build.gradle.kts or build.gradle not found under android/app")
}

func readKeyProperties(path string) (SigningConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SigningConfig{}, err
	}

	cfg := SigningConfig{}
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

	if err := validateSigningConfig(cfg); err != nil {
		return SigningConfig{}, err
	}
	return cfg, nil
}

func signingConfigFromCredentialsFile(path string) (SigningConfig, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return SigningConfig{}, fmt.Errorf("empty credentials path")
	}

	file, err := os.Open(path)
	if err != nil {
		return SigningConfig{}, err
	}
	defer file.Close()

	cfg := SigningConfig{}
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
		return SigningConfig{}, err
	}

	if cfg.StoreFile == "" {
		return SigningConfig{}, fmt.Errorf("credentials file missing keystore path")
	}
	if cfg.KeyPassword == "" {
		cfg.KeyPassword = cfg.StorePassword
	}
	if cfg.KeyAlias == "" {
		cfg.KeyAlias = defaultKeyAlias
	}

	if err := validateSigningConfig(cfg); err != nil {
		return SigningConfig{}, err
	}
	return cfg, nil
}

func validateSigningConfig(cfg SigningConfig) error {
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
