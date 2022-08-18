package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh cli-extension-demo",
	Short: "This extension prints out all the input user provides",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stage 1 of demo is done")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
