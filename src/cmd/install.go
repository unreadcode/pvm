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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "Install the specified PHP version.",
	Long: `You can quickly install a PHP version
Alias i
usage example:
	pvm install 8.0
	pvm i 8.0
`,
	Run: installRun,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		utils.PrintMsg("Please specify the PHP version to install.", "Warning", 1)
	}
	version := args[0]
	if utils.IsInstalled(version) {
		utils.PrintMsg(fmt.Sprintf("PHP v%s is already installed.", version), "Warning", 1)
	}
	install(version)
}

// 安装指定版本
func install(version string) {
	releases, err := utils.GetPHPReleases()
	if err != nil {
		utils.PrintMsg("Failed to get PHP releases.", "Error", 1)
	}

	zipInfo := releases[version]
	if zipInfo.Path == "" {
		utils.PrintMsg(fmt.Sprintf("PHP v%s is not found.", version), "Error", 1)
	}

	fileName, err := utils.Download(zipInfo.Path, version)
	if err != nil {
		utils.PrintMsg("Failed to download PHP release.", "Error", 1)
	}

	if err = utils.Unzip(fileName, version); err != nil {
		utils.PrintMsg("Failed to unzip PHP release.", "Error", 1)
	}

	if err = utils.CopyIni(version); err != nil {
		utils.PrintMsg("Failed to copy php.ini.", "Error", 1)
	}

	utils.PrintMsg(fmt.Sprintf("PHP v%s installed successfully.", version), "Success", 0)
}
