package pkgrename

import (
	"fmt"
	"strings"
)

func Run(projectPath string, answers []string) error {
	if len(answers) < 2 {
		return fmt.Errorf("update package name: missing old or new package name")
	}

	oldName := strings.TrimSpace(answers[0])
	newName := strings.TrimSpace(answers[1])
	if oldName == "" || newName == "" {
		return fmt.Errorf("update package name: old and new package names are required")
	}

	fmt.Printf("Updating package name in %s\n", projectPath)
	fmt.Printf("  old package: %s\n", oldName)
	fmt.Printf("  new package: %s\n", newName)

	result, err := Update(projectPath, oldName, newName)
	if err != nil {
		return err
	}

	fmt.Printf("  files updated: %d\n", result.FilesUpdated)
	fmt.Printf("  directories renamed: %d\n", result.DirsRenamed)
	for _, file := range result.UpdatedFiles {
		fmt.Printf("    updated: %s\n", file)
	}
	for _, dir := range result.RenamedDirs {
		fmt.Printf("    renamed: %s\n", dir)
	}

	if result.FilesUpdated == 0 && result.DirsRenamed == 0 {
		fmt.Println("  warning: no occurrences of the old package name were found")
	}

	return nil
}
