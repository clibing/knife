package system

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/clibing/knife/cmd/system/pkg"
	"github.com/spf13/cobra"
)

var App map[string]pkg.Application = make(map[string]pkg.Application)

// macOSCmd represents the cron command
var macOSCmd = &cobra.Command{
	Use:   "install",
	Short: "install 初始化一些软件",
	Long: `安装macoOS初始化

1. 主要应用于macos系统

需要先安装 Command Line Tools (CLT) 工具， 链接：https://developer.apple.com/download/all

01. 安装.vimrc规范文件
02. 安装homebrew
03. 安装golang
04. 安装iterm2终端

2. 安装golang, vim
knife system install -s golang -s vim

3. 强制配置vim
knife system install -o -s vim
`,

	Run: func(c *cobra.Command, arg []string) {
		overwrite, _ := c.Flags().GetBool("overwrite")
		log.Println("是否覆盖安装", overwrite)
		values, _ := c.Flags().GetStringSlice("select")

		if len(values) > 0 {
			log.Println("选择性执行安装: ", strings.Join(values, ","))
			for _, key := range values {
				if len(key) == 0 {
					continue
				}
				log.Println("正在安装:", key)
				execute(overwrite, App[key])
			}
			return
		}
		log.Println("执行默认安装")
		for _, app := range App {
			log.Println("正在安装:", app.Key())
			execute(overwrite, app)
		}
	},
}

func execute(overwrite bool, run pkg.Application) {
	pd := run.GetPackage()
	check := run.Before(pd, overwrite)
	if check {
		run.Install(pd)
	}
	run.Upgrade(pd)
	run.After(pd)

}

func init() {
	add([]pkg.Application{
		&pkg.Brew{},
		&pkg.Vim{},
		&pkg.Golang{},
		pkg.NewOhmyzshPlugin(),

		&pkg.GitflowControl{},

		&pkg.ITerm2{},
	},
	)

	help := []string{}
	for key := range App {
		help = append(help, key)
	}
	sort.Strings(help)

	macOSCmd.Flags().BoolP("overwrite", "o", false, "是否覆盖安装")
	macOSCmd.Flags().StringSliceP("select", "s", nil, fmt.Sprintf("可选择:[ %s ]", strings.Join(help, " | ")))
}

func NewMacOSCmd() *cobra.Command {
	return macOSCmd
}

func add(app []pkg.Application) {
	for _, value := range app {
		App[value.Key()] = value
	}
}
