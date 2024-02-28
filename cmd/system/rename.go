package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clibing/knife/cmd/debug"
	"github.com/clibing/knife/internal/utils"
	"github.com/spf13/cobra"
)

// rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "将指定目录下的文件，重新命名。",
	Long: `example: 
- ./
 |- 宝宝巴士-P1-你的名字.mp3

将当前目录下名字抹除重复部分 "宝宝巴士-" 为 ""
knife system rename -p ./ -s "宝宝巴士-" -t ""`,
	Run: func(c *cobra.Command, arg []string) {
		d := debug.NewDebug(c)
		p, _ := c.Flags().GetString("path")
		result := utils.Scan(d, p)

		source, _ := c.Flags().GetString("source")
		target, _ := c.Flags().GetString("target")

		exec, _ := c.Flags().GetBool("exec")

		for _, f := range result {
			parent := filepath.Dir(f)
			name := filepath.Base(f)
			newName := strings.ReplaceAll(name, source, target)

			result := filepath.Join(parent, newName)
			if exec {
				err := os.Rename(f, result)
				fmt.Printf("重命令: [%s] to [%s], result: %t\n", f, result, err == nil)
			} else {
				fmt.Printf("预检查: [%s] to [%s]\n", f, result)
			}
		}
	},
}

func init() {
	renameCmd.Flags().StringP("path", "p", "./", "扫描目录")
	renameCmd.Flags().StringP("filter", "f", "", "过滤部分")
	renameCmd.Flags().StringP("source", "s", "", "查找内容")
	renameCmd.Flags().StringP("target", "t", "", "目标内容")
	renameCmd.Flags().BoolP("exec", "e", false, "是否执行替换")
}

func NewRenameCmd() *cobra.Command {
	return renameCmd
}
