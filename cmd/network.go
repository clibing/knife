package cmd


import (
	"github.com/knife/cmd/network"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "net",
	Short: `网络处理器: ip, wol`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	networkCmd.AddCommand(network.NewIpCmd(), network.NewWolCmd())

	// 转换器
	rootCmd.AddCommand(networkCmd)
}