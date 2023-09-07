package cmd

import (
	"github.com/clibing/knife/cmd/common"
	"github.com/spf13/cobra"
)

var commonCmd = &cobra.Command{
	Use:   "common",
	Short: `通用工具: random, time`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	commonCmd.AddCommand(
		common.NewTimeCmd(),
		common.NewRandomCmd(),
	)

	// 通用工具类
	rootCmd.AddCommand(commonCmd)
}
