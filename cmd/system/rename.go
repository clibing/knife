package system

import (
	"fmt"

	"github.com/clibing/knife/cmd/debug"
	"github.com/clibing/knife/internal/utils"
	"github.com/spf13/cobra"
)

// rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename file",
	Long:  `详细使用: `,
	Run: func(c *cobra.Command, arg []string) {
		d := debug.NewDebug(c)
		p, _ := c.Flags().GetString("path")
		result := utils.Scan(d, p)

		t, _ := c.Flags().GetString("target")
		r, _ := c.Flags().GetString("replace")

		for _, name := range result {
			fmt.Println(t, r, name)
		}

	},
}

func init() {
	renameCmd.Flags().StringP("path", "p", "./", "扫描目录")
	renameCmd.Flags().StringP("filter", "f", "", "过滤部分")
	renameCmd.Flags().StringP("target", "t", "", "被替换的部分")
	renameCmd.Flags().StringP("replace", "r", "", "替换后的内容")
}

func NewRenameCmd() *cobra.Command {
	return renameCmd
}
