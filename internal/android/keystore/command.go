package keystore

import (
	"fmt"
	"path/filepath"

	"martini/internal/flutter/gradle"
	shared "martini/internal/keystore"
)

func Run(projectPath string, opts Options) error {
	if err := PromptMissing(&opts); err != nil {
		return err
	}

	cfg := opts.Config()
	if cfg.StorePassword == "" {
		return fmt.Errorf("store password is required")
	}

	fmt.Printf("Generating upload keystore for %s\n", projectPath)

	result, err := shared.Generate(projectPath, cfg)
	if err != nil {
		return err
	}

	keyPropertiesPath := filepath.Join(projectPath, "key.properties")
	if err := gradle.WriteKeyProperties(keyPropertiesPath, gradle.Properties{
		StoreFile:     result.FileName,
		StorePassword: result.StorePassword,
		KeyPassword:   result.KeyPassword,
		KeyAlias:      result.Alias,
	}); err != nil {
		return err
	}

	fmt.Printf("  home keystore: %s\n", result.HomePath)
	fmt.Printf("  project keystore: %s\n", result.ProjectPath)
	fmt.Printf("  credentials: %s\n", result.CredentialsPath)
	fmt.Printf("  key.properties: %s\n", keyPropertiesPath)
	fmt.Println("  note: build.gradle.kts was not modified")
	fmt.Println("  status: upload keystore generated")

	return nil
}
