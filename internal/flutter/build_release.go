package flutter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runBuildRelease(session *Session) error {
	projectPath := session.ProjectPath
	if err := validateFlutterProject(projectPath); err != nil {
		return err
	}

	fmt.Printf("Building release artifacts for %s\n", projectPath)

	if err := runFlutter(projectPath, "build", "apk", "--split-per-abi"); err != nil {
		return fmt.Errorf("flutter build apk: %w", err)
	}
	if err := runFlutter(projectPath, "build", "appbundle"); err != nil {
		return fmt.Errorf("flutter build appbundle: %w", err)
	}

	artifacts, err := collectReleaseArtifacts(projectPath)
	if err != nil {
		return err
	}
	if len(artifacts) == 0 {
		return fmt.Errorf("no release artifacts found under build/app/outputs")
	}

	zipPath, err := createReleaseZip(projectPath, artifacts)
	if err != nil {
		return err
	}

	fmt.Println("  release artifacts:")
	for _, artifact := range artifacts {
		fmt.Printf("    %s (%s)\n", artifact.Name, artifact.SHA1)
	}
	fmt.Printf("  release zip: %s\n", zipPath)
	fmt.Println("  status: release build completed")

	return nil
}

func validateFlutterProject(projectPath string) error {
	info, err := os.Stat(projectPath)
	if err != nil {
		return fmt.Errorf("project path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", projectPath)
	}
	pubspec := filepath.Join(projectPath, "pubspec.yaml")
	if _, err := os.Stat(pubspec); err != nil {
		return fmt.Errorf("pubspec.yaml not found in project root")
	}
	return nil
}

func runFlutter(projectPath string, args ...string) error {
	cmd := exec.Command("flutter", args...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}
