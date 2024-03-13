package utils

import (
	"io/fs"
	"path/filepath"

	"github.com/clibing/knife/cmd/debug"
)

/**
 * 扫描指定目录
 */

func Scan(debug *debug.Debug, path string) (result []string) {
	scanDir, _ := filepath.Abs(path)
	debug.ShowSame("🟣 scan path: %s", scanDir)

	// 扫描
	filepath.Walk(scanDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只保存非目录
		if !info.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	return
}
