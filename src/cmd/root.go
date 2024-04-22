//Package cmd
/*
Copyright © 2024 UnreadCode <i@unreadcode.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pvm",
	Short: "PHP version manager for Windows.",
	Long: `A PHP version management tool for Windows, allowing you to easily switch between your PHP versions!
Copyright (c) 2024 UnreadCode <i@unreadcode.com>
Visit https://github.com/unreadcode/pvm for more information.
`,
	Run: rootRun,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print PVM version number.")
}

func rootRun(cmd *cobra.Command, args []string) {
	version, _ := cmd.Flags().GetBool("version")
	if version {

	}
	cmd.Help()
}
