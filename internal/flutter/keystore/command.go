package keystore

import (
	"fmt"

	"martini/internal/flutter/gradle"
)

func Run(projectPath string, answers []string) error {
	cfg := configFromAnswers(answers)
	if cfg.StorePassword == "" {
		return fmt.Errorf("generate upload keystore: store password is required")
	}
	if cfg.KeyPassword == "" {
		cfg.KeyPassword = cfg.StorePassword
	}

	fmt.Printf("Generating upload keystore for %s\n", projectPath)

	result, err := Generate(projectPath, cfg)
	if err != nil {
		return err
	}

	if err := gradle.Apply(projectPath, gradle.Properties{
		StoreFile:     result.HomePath,
		StorePassword: result.StorePassword,
		KeyPassword:   result.KeyPassword,
		KeyAlias:      result.Alias,
	}); err != nil {
		return err
	}

	fmt.Printf("  home keystore: %s\n", result.HomePath)
	fmt.Printf("  project keystore: %s\n", result.ProjectPath)
	fmt.Printf("  credentials: %s\n", result.CredentialsPath)
	fmt.Printf("  key.properties: %s/android/key.properties\n", projectPath)
	fmt.Println("  build.gradle.kts updated for release signing")
	fmt.Println("  status: upload keystore setup completed")

	return nil
}
