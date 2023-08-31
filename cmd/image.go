package cmd

import (
	"github.com/knife/cmd/image"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: `图片处理器: qrcode, base64 image`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	// 增加二维码处理器
	imageCmd.AddCommand(image.NewQrcodeCmd())

	// 转换器
	rootCmd.AddCommand(imageCmd)
}
