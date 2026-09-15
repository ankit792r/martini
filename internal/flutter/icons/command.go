package icons

import (
	"fmt"
	"strings"
)

func Run(projectPath string, answers []string) error {
	if len(answers) == 0 || strings.TrimSpace(answers[0]) == "" {
		return fmt.Errorf("update icons: IconKitchen output path is required")
	}
	iconsPath := strings.TrimSpace(answers[0])

	resolved, err := ResolveIconsPath(iconsPath)
	if err != nil {
		return err
	}

	fmt.Printf("Updating icons in %s\n", projectPath)
	fmt.Printf("  icon source: %s\n", resolved)

	result, err := Update(projectPath, iconsPath)
	if err != nil {
		return err
	}

	fmt.Printf("  android files: %d\n", result.AndroidFiles)
	fmt.Printf("  ios files: %d\n", result.IOSFiles)
	fmt.Printf("  web files: %d\n", result.WebFiles)
	fmt.Println("  status: icons updated")

	return nil
}
