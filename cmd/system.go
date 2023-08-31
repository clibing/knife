package cmd

import (
	"github.com/knife/cmd/system"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:     "system",
	Aliases: []string{"sys"},
	Short:   `系统工具: arch, monitor, upgrade`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	systemCmd.AddCommand(system.NewArchCmd(), system.NewMonitorCmd(), system.NewUpgradeCmd())

	// 转换器
	rootCmd.AddCommand(systemCmd)
}
