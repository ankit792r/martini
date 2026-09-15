package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var flutterCmd = &cobra.Command{
	Use:   "flutter",
	Short: "Flutter related commands",
	Long:  "Select tasks to perform on a specific Flutter project.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flutter command")
	},
}

func init() {
	rootCmd.AddCommand(flutterCmd)
}
