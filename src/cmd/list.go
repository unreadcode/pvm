//Package cmd
/*
Copyright © 2024 UnreadCode <i@unreadcode.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"pvm/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the installed PHP versions.",
	Long: `List the installed PHP versions.
Alias ls
Add the -available or - a parameter to list all Available  PHP versions.
usage example:
	pvm list
	pvm ls
	pvm list -a
	pvm ls -a
`,
	Run: listRun,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("available", "a", false, "List all Available PHP versions.")
}

func listRun(cmd *cobra.Command, _ []string) {
	available, _ := cmd.Flags().GetBool("available")
	if available {
		listAvailable()
	} else {
		listInstalled()
	}
}

// 列出所有可用的PHP版本
func listAvailable() {
	versions, err := utils.GetPHPReleases()
	if err != nil {
		utils.PrintMsg(err.Error(), "Error", 1)
	}
	if len(versions) <= 0 {
		utils.PrintMsg("no available PHP versions found", "Warning", 1)
	}
	fmt.Println("Available PHP versions:")
	for v := range versions {
		utils.PrintMsg(fmt.Sprintf("      v%s", v), "Info", 888)
	}
	utils.PrintMsg(fmt.Sprintf("\nThis list from %s", utils.RELEASES), "Warning", 888)
}

// 列出已安装的PHP版本
func listInstalled() {
	installed, err := utils.GetInstalledVersions()
	if err != nil {
		utils.PrintMsg(err.Error(), "Error", 1)
	}
	if len(installed) <= 0 {
		utils.PrintMsg("no installed PHP versions found", "Warning", 1)
	}
	currentVersion := utils.GetCurrentPhpVersion()
	utils.PrintMsg("Installed PHP versions:", "Info", 888)
	for _, v := range installed {
		if v == fmt.Sprintf("v%s", currentVersion) {
			utils.PrintMsg(fmt.Sprintf("    * %s", v), "Success", 888)
			continue
		}
		utils.PrintMsg(fmt.Sprintf("      %s", v), "Info", 888)
	}

}
