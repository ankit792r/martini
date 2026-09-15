package cmd

import (
	"fmt"
	"os"

	"martini/internal/devbox"

	"github.com/spf13/cobra"
)

var devboxCmd = &cobra.Command{
	Use:   "devbox",
	Short: "Devbox related commands",
	Long:  "Add and manage devbox configuration in selected projects.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := devbox.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devboxCmd)
}
