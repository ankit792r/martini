package cmd

import (
	"fmt"
	"os"

	"martini/internal/android"
	androidkeystore "martini/internal/android/keystore"

	"github.com/spf13/cobra"
)

var (
	androidProjectPath        string
	androidGenerateKeystore   bool
	androidKeystoreName       string
	androidStorePassword      string
	androidKeyPassword        string
	androidKeyAlias           string
	androidCommonName         string
)

var androidCmd = &cobra.Command{
	Use:   "android",
	Short: "Android Studio project commands",
	Long:  "Tools for native Android projects created with Android Studio.",
	Run: func(cmd *cobra.Command, args []string) {
		if !androidGenerateKeystore {
			cmd.Help()
			return
		}

		projectPath, err := android.ResolveProjectPath(androidProjectPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		err = androidkeystore.Run(projectPath, androidkeystore.Options{
			KeystoreName:  androidKeystoreName,
			StorePassword: androidStorePassword,
			KeyPassword:   androidKeyPassword,
			KeyAlias:      androidKeyAlias,
			CommonName:    androidCommonName,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(androidCmd)

	androidCmd.Flags().StringVarP(&androidProjectPath, "path", "p", "", "Android project path (default: current directory)")
	androidCmd.Flags().BoolVar(&androidGenerateKeystore, "generate-upload-keystore", false, "Generate an upload keystore and key.properties")
	androidCmd.Flags().StringVar(&androidKeystoreName, "keystore-name", "", "Keystore file name without .jks extension")
	androidCmd.Flags().StringVar(&androidStorePassword, "store-password", "", "Keystore store password")
	androidCmd.Flags().StringVar(&androidKeyPassword, "key-password", "", "Key password (default: store password)")
	androidCmd.Flags().StringVar(&androidKeyAlias, "key-alias", "", "Key alias")
	androidCmd.Flags().StringVar(&androidCommonName, "common-name", "", "Certificate common name")
}
