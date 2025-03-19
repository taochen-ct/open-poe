package path

import (
	"os"
)

// WorkPath 当前工作路径
func WorkPath() string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}
