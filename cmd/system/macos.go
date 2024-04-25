package system

import (
	"github.com/clibing/knife/cmd/system/pkg"
	"github.com/spf13/cobra"
)

var App []pkg.Application

// macOSCmd represents the cron command
var macOSCmd = &cobra.Command{
	Use:   "macos",
	Short: "macos 初始化",
	Long: `安装macoOS初始化

需要先安装 Command Line Tools (CLT) 工具， 链接：https://developer.apple.com/download/all

01. 安装.vimrc规范文件

.`,

	Run: func(c *cobra.Command, arg []string) {
		overwrite, _ := c.Flags().GetBool("overwrite")

		for _, app := range App {
			pd := app.GetPackage()
			check := app.Before(pd, overwrite)
			if check {
				app.Install(pd)
			}
			app.Upgrade(pd)
			app.After(pd)
		}
	},
}

func init() {
	macOSCmd.Flags().BoolP("overwrite", "o", false, "是否覆盖安装")

	App = append(App,
		pkg.NewOhmyzsh(),
		// &pkg.Brew{},
		// &pkg.Vim{}, // pass
		// &pkg.Go{},

		// application
		// &pkg.ITerm2{},
	)
}

func NewMacOSCmd() *cobra.Command {
	return macOSCmd
}
