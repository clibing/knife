package system

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

// rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename file",
	Long:  `详细使用: `,
	Run: func(_ *cobra.Command, arg []string) {
		
		
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
