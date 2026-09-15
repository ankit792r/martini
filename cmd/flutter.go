package cmd

import (
	"fmt"
	"os"

	"martini/internal/flutter"

	"github.com/spf13/cobra"
)

var (
	flutterProjectPath         string
	flutterUpdatePackageName   bool
	flutterGenerateKeystore    bool
	flutterUpdateSigningConfig bool
)

var flutterCmd = &cobra.Command{
	Use:   "flutter",
	Short: "Flutter related commands",
	Long:  "Select tasks to perform on a specific Flutter project.",
	Run: func(cmd *cobra.Command, args []string) {
		flags := map[string]bool{
			"update-package-name":      flutterUpdatePackageName,
			"generate-upload-keystore": flutterGenerateKeystore,
			"update-signing-config":    flutterUpdateSigningConfig,
		}

		preselected := flutter.CommandsFromFlags(flags)
		if err := flutter.Run(flutter.Options{
			ProjectPath: flutterProjectPath,
			Preselected: preselected,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(flutterCmd)

	flutterCmd.Flags().StringVarP(&flutterProjectPath, "path", "p", "", "Flutter project path")
	flutterCmd.Flags().BoolVar(&flutterUpdatePackageName, "update-package-name", false, "Update Gradle package name")
	flutterCmd.Flags().BoolVar(&flutterGenerateKeystore, "generate-upload-keystore", false, "Generate upload keystore")
	flutterCmd.Flags().BoolVar(&flutterUpdateSigningConfig, "update-signing-config", false, "Update signing config")
}
