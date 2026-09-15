package flutter

import (
	"fmt"
)

func runGenerateUploadKeystore(session *Session) error {
	cfg := keystoreConfigFromAnswers(session.AnswersFor(CmdGenerateUploadKeystore))
	if cfg.StorePassword == "" {
		return fmt.Errorf("generate upload keystore: store password is required")
	}
	if cfg.KeyPassword == "" {
		cfg.KeyPassword = cfg.StorePassword
	}

	fmt.Printf("Generating upload keystore for %s\n", session.ProjectPath)

	result, err := generateUploadKeystore(session.ProjectPath, cfg)
	if err != nil {
		return err
	}

	if err := configureAndroidSigning(session.ProjectPath, result); err != nil {
		return err
	}

	fmt.Printf("  home keystore: %s\n", result.HomePath)
	fmt.Printf("  project keystore: %s\n", result.ProjectPath)
	fmt.Printf("  credentials: %s\n", result.CredentialsPath)
	fmt.Printf("  key.properties: %s/android/key.properties\n", session.ProjectPath)
	fmt.Println("  build.gradle.kts updated for release signing")
	fmt.Println("  status: upload keystore setup completed")

	return nil
}
