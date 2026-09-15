package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devboxCmd = &cobra.Command{
	Use:   "devbox",
	Short: "Devbox related commands",
	Long:  "Add and manage devbox configuration in selected projects.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("devbox command")
	},
}

func init() {
	rootCmd.AddCommand(devboxCmd)
}
