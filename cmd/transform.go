package cmd

import (
	"github.com/clibing/knife/cmd/transform"
	"github.com/spf13/cobra"
)

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: `转换器: json, markdown, url`,
	Long: `转换器

1. markdown转换器(html ←→ md)
2. url encoding, url decoding`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	// markdown 转换
	mdCmd := transform.NewMarkdown()
	transformCmd.AddCommand(mdCmd)

	// url encoding decoding
	urlCmd := transform.NewUrlEncoding()
	transformCmd.AddCommand(urlCmd)

	// add text convert
	// transformCmd.AddCommand(transform.NewJsonConvert())

	// 转换器
	rootCmd.AddCommand(transformCmd)
}
