package utils

import (
	"fmt"
	"strconv"
)

const (
	B  = 1024
	KB = B * 1024
	MB = KB * 1024
	GB = KB * 1024
	TB = GB * 1024
	PB = TB * 1024
)

/**
 * 美化格式化 转 单位 和 值
 */
func BeautifyUnit(size int64) (unit, value string) {
	unit, value = baseBeautifySize(size)
	return
}

/**
 * 美化格式化 转 字符串 (值 单位)
 * 2048B -> 2 K
 */
func BeautifyValue(size int64) string {
	unit, value := baseBeautifySize(size)
	return fmt.Sprintf("%s %s", value, unit)
}

/**
 * 美化 文件的大小
 */
func baseBeautifySize(fileSize int64) (unit, value string) {
	if fileSize < B {
		return "B", fmt.Sprintf("%.2f", float64(fileSize)/float64(1))
	} else if fileSize < KB {
		return "KB", fmt.Sprintf("%.2f", float64(fileSize)/float64(B))
	} else if fileSize < MB {
		return "MB", fmt.Sprintf("%.2f", float64(fileSize)/float64(KB))
	} else if fileSize < GB {
		return "GB", fmt.Sprintf("%.2f", float64(fileSize)/float64(MB))
	} else if fileSize < TB {
		return "TB", fmt.Sprintf("%.2f", float64(fileSize)/float64(GB))
	} else if fileSize < PB {
		return "PB", fmt.Sprintf("%.2f", float64(fileSize)/float64(TB))
	} else {
		return "", strconv.FormatInt(fileSize, 10)
	}
}
