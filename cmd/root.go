package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "martini",
	Short:   "Post development multipurpose tool",
	Long:    "Martini is a tool that helps you manage your development tasks",
	Example: "  martini devbox\n  martini flutter",
	// SilenceUsage:      true,
	// SilenceErrors:     true,
	// DisableAutoGenTag: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
