//Package cmd
/*
Copyright © 2024 UnreadCode <i@unreadcode.com>

*/

package cmd

import (
	"github.com/spf13/cobra"
	"pvm/utils"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Aliases: []string{"uni"},
	Short:   "Uninstall a PHP version",
	Long: `Uninstall a PHP version from the system.
Alias uni
usage example:
	pvm uninstall 8.0
	pvm uni 8.0
`,
	Run: uninstallRun,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func uninstallRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		utils.PrintMsg("Please specify the PHP version to uninstall.", "Warning", 1)
	}
	version := args[0]
	// 检查版本是否为当前PHP版本
	if version == utils.GetCurrentPhpVersion() {
		utils.PrintMsg("You cannot uninstall the current PHP version.", "Warning", 1)
	}
	if err := utils.UninstallPhpVersion(version); err != nil {
		utils.PrintMsg(err.Error(), "Error", 1)
	}
	utils.PrintMsg("PHP version "+version+" uninstalled successfully.", "Success", 0)
}
