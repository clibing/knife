package system

import (
	"github.com/clibing/knife/cmd/system/pkg"
	"github.com/spf13/cobra"
)

// ohmyzshCmd represents the cron command
var ohmyzshCmd = &cobra.Command{
	Use:   "ohmyzsh",
	Short: "安装ohmyzsh初始化",
	Long: `安装macoOS初始化

注意：内部调用官忘sh脚本安装，安装后会exit 1

需要先安装 Command Line Tools (CLT) 工具， 链接：https://developer.apple.com/download/all.`,

	Run: func(c *cobra.Command, arg []string) {
		overwrite, _ := c.Flags().GetBool("overwrite")
		ozh := pkg.NewOhmyzsh()

		pd := ozh.GetPackage()
		check := ozh.Before(pd, overwrite)
		if check {
			ozh.Install(pd)
		}
		ozh.Upgrade(pd)
		ozh.After(pd)
	},
}

func init() {
	ohmyzshCmd.Flags().BoolP("overwrite", "o", false, "是否覆盖安装")
}

func NewOhmyzshCmd() *cobra.Command {
	return ohmyzshCmd
}
