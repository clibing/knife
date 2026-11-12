package utils

import (
	"fmt"
)

const (
	B  = 1
	KB = B * 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
)

/**
 * 美化格式化 转 单位 和 值
 */
func BeautifyUnit(size int64) (unit, value string) {
	return baseBeautifySize(size)
}

/**
 * 美化格式化 转 字符串 (值 单位)
 * 2048B -> 2.00 KB
 */
func BeautifyValue(size int64) string {
	unit, value := baseBeautifySize(size)
	return fmt.Sprintf("%s %s", value, unit)
}

func baseBeautifySize(fileSize int64) (unit, value string) {
	switch {
	case fileSize < KB:
		return "B", fmt.Sprintf("%.2f", float64(fileSize)/B)
	case fileSize < MB:
		return "KB", fmt.Sprintf("%.2f", float64(fileSize)/KB)
	case fileSize < GB:
		return "MB", fmt.Sprintf("%.2f", float64(fileSize)/MB)
	case fileSize < TB:
		return "GB", fmt.Sprintf("%.2f", float64(fileSize)/GB)
	case fileSize < PB:
		return "TB", fmt.Sprintf("%.2f", float64(fileSize)/TB)
	default:
		return "PB", fmt.Sprintf("%.2f", float64(fileSize)/PB)
	}
}
