package system

import (
	"github.com/spf13/cobra"
)

// beautifyCmd represents the cpu command
var beautifyCmd = &cobra.Command{
	Use:   "beautify",
	Short: "字节美化工具, 反向计算字节",
	Run: func(c *cobra.Command, _ []string) {
		c.Help()
	},
}

func init() {
	beautifyCmd.Flags().StringP("byte", "b", "", "")
}

func NewBeautifyCmd() *cobra.Command {
	return beautifyCmd
}
