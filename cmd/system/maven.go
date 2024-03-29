package system

import (
	"fmt"
	"os"

	"github.com/clibing/knife/cmd/debug"
	"github.com/clibing/knife/internal/maven"
	"github.com/clibing/knife/internal/utils"
	"github.com/spf13/cobra"
)

// mvnCmd represents the arch command
var mvnCmd = &cobra.Command{
	Use:     "maven",
	Aliases: []string{"mvn"},
	Short:   "Clean Maven snapshot dependencies",
	Example: `
1. 扫描并分析当前目录下的jar包
knife system maven ./	

2. 扫描、分析并删除
knife system maven -c ./`,

	Run: func(c *cobra.Command, args []string) {
		debug := debug.NewDebug(c)

		if len(args) == 0 {
			debug.Debug("暂无扫描的path")
			c.Help()
			return
		}
		clear, _ := c.Flags().GetBool("clear")
		var freeUpspace int64
		debug.ShowSame("🔵 clean maven starting.")
		has := false
		for _, p := range args {
			data := maven.Doing(debug, p)
			for _, v := range data {
				for _, w := range v.Snapshot {
					has = true
					if w.Deleted {
						info, err := os.Stat(w.FullName)
						if err != nil {
							debug.Debug("remove to failed, err: %s, file: %s", err.Error(), w.FullName)
						} else {
							freeUpspace = freeUpspace + info.Size()
						}
						if clear {
							err = os.Remove(w.FullName)
							if err != nil {
								debug.Debug("remove to failed, err: %s, file: %s", err.Error(), w.FullName)
							}
							sha1 := fmt.Sprintf("%s.sha1", w.FullName)
							_, err = os.Stat(sha1)
							if err == nil || os.IsExist(err) {
								err = os.Remove(sha1)
								if err != nil {
									debug.Debug("remove to failed, err: %s, file: %s", err.Error(), w.FullName)
								}
							}
							debug.ShowSame("✅ successfully removed: %s", w.FullName)
						} else {
							debug.ShowSame("❓ preview automatically deleted when settings are deleted: %s", w.FullName)
						}
					} else {
						debug.ShowSame("🟡 skip delete: %s", w.FullName)
					}
				}
			}
		}
		if !has {
			debug.ShowSame("🟡 暂无需要处理的文件")
		}
		if freeUpspace > 0 {
			if !clear {
				debug.ShowSame("❗ will free disk: %s", utils.BeautifyValue(freeUpspace))
			} else {
				debug.ShowSame("🟤 free disk: %s", utils.BeautifyValue(freeUpspace))
			}
		}
		debug.ShowSame("🔵 clean maven end")
	},
}

func init() {
	mvnCmd.Flags().BoolP("clear", "c", false, "是否确认清理过期的snapshot jar")
}

func NewMavenCmd() *cobra.Command {
	return mvnCmd
}
