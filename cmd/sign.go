package cmd

import (
	"github.com/clibing/knife/cmd/sign"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: `签名密钥相关: rsa, md5, base64, sha1, sha128, sha256, sha512; 二次验证 otp`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	signCmd.AddCommand(
		sign.NewRsaCmd(),
		sign.NewMd5Cmd(),
		sign.NewShaCmd(),
		sign.NewBase64Cmd(),
		sign.NewAesCmd(),
		sign.NewOtpCmd(),
	)

	// 转换器
	rootCmd.AddCommand(signCmd)
}
