package signing

import (
	"fmt"
)

func Run(projectPath string) error {
	fmt.Printf("Updating signing config for %s\n", projectPath)

	cfg, err := Resolve(projectPath)
	if err != nil {
		return err
	}

	if err := Apply(projectPath, cfg); err != nil {
		return err
	}

	fmt.Printf("  key.properties: %s/android/key.properties\n", projectPath)
	fmt.Printf("  store file: %s\n", cfg.StoreFile)
	fmt.Println("  build.gradle updated for release signing")
	fmt.Println("  status: signing config updated")

	return nil
}
