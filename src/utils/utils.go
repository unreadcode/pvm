//Package utils
/*
Copyright © 2024 UnreadCode <i@unreadcode.com>
*/

package utils

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type PHPReleases map[string]Zip

type Zip struct {
	Path   string
	Size   string
	SHA256 string
}

var PvmRoot = os.Getenv(PVM_ROOT)

var MsgTypeMap = map[string]color.Attribute{
	"Error":   color.FgRed,
	"Warning": color.FgYellow,
	"Success": color.FgGreen,
	"Info":    color.FgCyan,
}

// PrintMsg 打印信息
func PrintMsg(message string, msgType string, code int) {
	if _, ok := MsgTypeMap[msgType]; !ok {
		msgType = "Info"
	}
	color.Set(MsgTypeMap[msgType])
	if code != 888 {
		fmt.Printf("[%s] %s\n", msgType, message)
	} else {
		fmt.Printf("%s\n", message)
	}
	color.Unset()
	if code != 888 {
		os.Exit(code)
	}
}

// GetInstalledVersions 获取已安装的版本
func GetInstalledVersions() ([]string, error) {
	phpDir, err := os.ReadDir(PvmRoot)
	// 读取目录失败
	if err != nil {
		return nil, err
	}
	var installed []string
	for _, f := range phpDir {
		// 跳过非文件夹
		if !f.IsDir() {
			continue
		}
		// 是有效的php版本
		if isValidPhpVersion(f.Name()) {
			installed = append(installed, f.Name())
		}
	}
	return installed, nil
}

// GetCurrentPhpVersion 获取当前的php版本
func GetCurrentPhpVersion() string {
	// 执行php -v命令
	cmd := exec.Command("php", "-v")
	outputBytes, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	outputStr := string(outputBytes)
	versionLine := strings.SplitN(outputStr, "\n", 1)[0]
	return parsePHPVersion(versionLine)
}

// GetPHPReleases 获取PHP维护中的所有发行版
func GetPHPReleases() (PHPReleases, error) {
	resp, err := http.Get(RELEASES)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &releases); err != nil {
		return nil, err
	}

	pattern := `^nts.*` + GetArch()
	var result = PHPReleases{}

	for versionNumber, val := range releases {
		for arch, zipPath := range val.(map[string]interface{}) {
			// 如果不是nts开头(非线程安全)或者不是当前系统架构结尾，则跳过
			if !regexp.MustCompile(pattern).MatchString(arch) {
				continue
			}
			zipPath := zipPath.(map[string]interface{})["zip"]
			zipInfo := zipPath.(map[string]interface{})
			version := Zip{}
			version.Path = fmt.Sprintf("%s%s", DOWNLOAD, zipInfo["path"].(string))
			version.Size = zipInfo["size"].(string)
			version.SHA256 = zipInfo["sha256"].(string)
			result[versionNumber] = version
		}
	}

	return result, nil
}

// GetArch 获取当前系统的架构
func GetArch() string {
	arch := runtime.GOARCH
	switch arch {
	case "386":
		arch = "x86"
	case "amd64":
		arch = "x64"
	default:
		arch = "x64"
	}
	return arch
}

// 判断给定的路径是否表示有效的PHP版本
func isValidPhpVersion(path string) bool {
	dirName := filepath.Base(path)
	// 如果不是v开头，也不是数字点分隔的，则不是有效的PHP版本
	if !regexp.MustCompile(`^v\d+\.\d+$`).MatchString(dirName) {
		return false
	}
	phpExe := filepath.Join(PvmRoot, path, "php.exe")
	fileInfo, err := os.Stat(phpExe)
	// 如果php.exe不存在，则不是有效的PHP版本
	if err != nil {
		return false
	}
	// 文件是否可执行
	return fileInfo.Mode().IsRegular()
}

// 从给定的行中解析出PHP版本
func parsePHPVersion(line string) string {
	pattern := `PHP (\d+\.\d+)`
	re := regexp.MustCompile(pattern)
	if match := re.FindStringSubmatch(line); match != nil {
		return match[1]
	}
	return "unknown"
}
