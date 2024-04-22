//Package utils
/*
Copyright © 2024 UnreadCode <i@unreadcode.com>
*/

package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

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
