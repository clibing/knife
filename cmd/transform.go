package cmd

import (
	"github.com/knife/cmd/transform"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "转换器, markdown转换器(html ←→ md)",
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	// markdown 转换
	mdCmd := transform.NewMarkdown()
	convertCmd.AddCommand(mdCmd)

	// 转换器
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mdCmd.PersistentFlags().String("foo", "", "A help for foo")
	// mdCmd.Flags().StringVarP(&source, "source", "s", "", "源文件")
	// mdCmd.Flags().StringVarP(&target, "target", "t", "", "目标文件")
	// mdCmd.Flags().BoolVarP(&direct, "direct", "d", false, "转换的目标类型，from->to，默认是HTML到markdown")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
