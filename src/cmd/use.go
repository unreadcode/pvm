//Package cmd
/*
Copyright © 2024 UnreadCode <i@unreadcode.com>

*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"pvm/utils"
	"regexp"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Switch to the specified PHP version",
	Long: `Quickly switch between different PHP versions.
usage example:
	pvm use 8.0
`,

	Run: useRun,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func useRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		utils.PrintMsg("Please specify a PHP version", "Warning", 1)
	}
	version := args[0]
	// 是一个有效的版本号
	if !regexp.MustCompile(`^\d+\.\d+$`).MatchString(version) {
		utils.PrintMsg("Invalid PHP version number.", "Error", 1)
	}
	// 是否已经在使用
	if currentVersion := utils.GetCurrentPhpVersion(); currentVersion == version {
		utils.PrintMsg(fmt.Sprintf("Already using PHP v%s", version), "Info", 0)
	}
	// 是否已经安装
	if !utils.IsInstalled(version) {
		utils.PrintMsg(fmt.Sprintf("PHP v%s is not installed.", version), "Error", 1)
	}
	// 切换到指定版本
	if err := utils.SwitchToVersion(version); err != nil {
		utils.PrintMsg(err.Error(), "Error", 1)
	}
	utils.PrintMsg(fmt.Sprintf("Switched to PHP v%s", version), "Success", 0)
}
