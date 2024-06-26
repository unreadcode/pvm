//Package utils
/*
Copyright © 2024 UnreadCode <i@unreadcode.com>
*/

package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
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
var PhpPath = os.Getenv(PHP_PATH)

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

// IsInstalled 给定的版本是否已安装
func IsInstalled(version string) bool {
	installed, _ := GetInstalledVersions()
	for _, v := range installed {
		if v == fmt.Sprintf("v%s", version) {
			return true
		}
	}
	return false
}

// Download 下载PHP版本
func Download(downloadUrl string, version string) (string, error) {
	if _, err := os.Stat(PvmRoot); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(PvmRoot, os.ModePerm); err != nil {
			return "", err
		}
	}
	PrintMsg(fmt.Sprintf("Downloading PHP v%s...", version), "Info", 888)
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	tempDir := os.TempDir()
	fileName := fmt.Sprintf("download_php_v%s.zip", version)
	tmpFile, err := os.CreateTemp(tempDir, fileName)
	if err != nil {
		return "", err
	}

	defer func(tmpFile *os.File) {
		err := tmpFile.Close()
		if err != nil {
			return
		}
	}(tmpFile)

	color.Set(color.FgCyan)
	// 创建进度条
	bar := pb.Full.Start64(resp.ContentLength)
	bar.Set(pb.Bytes, true)
	bar.Start()
	reader := bar.NewProxyReader(resp.Body)
	if _, err = io.Copy(tmpFile, reader); err != nil {
		return "", err
	}
	bar.Finish()
	color.Unset()

	return tmpFile.Name(), nil
}

// Unzip 解压PHP版本
func Unzip(zipfilePath string, version string) error {
	archive, err := zip.OpenReader(zipfilePath)
	if err != nil {
		return err
	}
	defer func(archive *zip.ReadCloser) {
		err := archive.Close()
		if err != nil {
			return
		}
	}(archive)
	// 解压目标目录
	targetDir := filepath.Join(PvmRoot, "v"+version)
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("PHP v%s already exists", version)
	}

	for _, item := range archive.File {
		filePath := filepath.Join(targetDir, item.Name)

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		if item.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, item.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := item.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, fileInArchive); err != nil {
			return err
		}
	}
	return nil
}

// CopyIni 复制ini文件
func CopyIni(version string) error {
	iniFilePath := filepath.Join(PvmRoot, "v"+version, "php.ini")
	productionIniFilePath, err := os.ReadFile(iniFilePath + "-production") //php.ini-production
	if err != nil {
		return err
	}
	if err := os.WriteFile(iniFilePath, productionIniFilePath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// UninstallPhpVersion 卸载PHP版本
func UninstallPhpVersion(version string) error {
	//版本是否存在
	if _, err := os.Stat(filepath.Join(PvmRoot, "v"+version)); err != nil {
		return fmt.Errorf("PHP v%s does not exist", version)
	}
	// 卸载版本
	if err := os.RemoveAll(filepath.Join(PvmRoot, "v"+version)); err != nil {
		return fmt.Errorf("failed to uninstall PHP v%s: %v", version, err)
	}
	return nil
}

// SwitchToVersion 切换到指定版本
func SwitchToVersion(version string) error {
	targetDir := filepath.Join(PvmRoot, fmt.Sprintf("v%s", version))
	if hasPermission() {
		os.Remove(PhpPath)
		if err := os.Symlink(targetDir, PhpPath); err != nil {
			return err
		}
		return nil
	} else {
		uacRun("pvm", "use", version)
	}
	return nil
}

// 是否具有管理员权限
func hasPermission() bool {
	if _, err := os.Open("\\\\.\\PHYSICALDRIVE0"); err != nil {
		return false
	}
	return true
}

// 请求管理员权限运行程序
func uacRun(command ...string) {
	//This executable file comes from https://jianqiezhushou.com/
	//Original author https://jpassing.com/2007/12/08/launch-elevated-processes-from-the-command-line/
	cmd := exec.Command("elevate.exe", "-k")
	cmd.Args = append(cmd.Args, command...)
	cmd.Run()
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
