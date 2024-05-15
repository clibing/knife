package cmd

import (
	"github.com/clibing/knife/cmd/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: `服务器相关: static, webdav`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	// 增加二维码处理器
	serverCmd.AddCommand(server.NewFileServer(), server.NewWebdavServer())

	// 转换器
	rootCmd.AddCommand(serverCmd)
}
