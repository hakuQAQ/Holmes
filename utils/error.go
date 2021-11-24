package utils

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
)

// 输出错误并退出
func OptionsError(message string, exitCode int) {
	fmt.Println("[-] " + message)
	cli.Exit(2)
}
