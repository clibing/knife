package sign

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: `OTP 生成验证码`,
	Long: `说明
通用配置: 刷新30秒， Digits: DigitsSix, Algorithm: SHA1

1. 解析二维码内容
knife image qrcode -d -i source.png

2. 生成验证码

* 应用：knife
* 账号: knife@linuxcrypt.cn
* 秘钥: ABC
knife sign otp -i knife -a knife@linuxcrypt.cn -s ABC

`,
	Run: func(c *cobra.Command, args []string) {
		qr, _ := c.Flags().GetBool("qr")
		account, _ := c.Flags().GetString("account")
		issuer, _ := c.Flags().GetString("issuer")
		secret, _ := c.Flags().GetString("secret")

		period, _ := c.Flags().GetInt("period")
		opts := totp.GenerateOpts{
			Issuer:      issuer,
			AccountName: account,
			Secret:      []byte(secret),
		}
		if period > 0 {
			opts.Period = uint(period)
		}

		// 生成秘钥
		key, err := totp.Generate(opts)

		if err != nil {
			fmt.Println("生成秘钥异常:", err)
			return
		}

		// 生成验证码
		passcode, err := totp.GenerateCode(key.Secret(), time.Now())
		if err != nil {
			fmt.Println("生成验证码失败:", err)
			return
		}
		fmt.Println("验证码:", passcode)

		// 生成二维码
		if qr {
			output, _ := c.Flags().GetString("output")
			width, _ := c.Flags().GetInt("width")
			height, _ := c.Flags().GetInt("height")

			var buf bytes.Buffer
			img, err := key.Image(width, height)
			if err != nil {
				fmt.Println("输出验证码失败:", err)
				return
			}

			png.Encode(&buf, img)
			os.WriteFile(output, buf.Bytes(), 0644)
		}
	},
}

func init() {
	otpCmd.Flags().StringP("secret", "s", "", "OTP 验证秘钥")
	otpCmd.Flags().StringP("issuer", "i", "knife", "应用名字")
	otpCmd.Flags().StringP("account", "a", "knife@linuxcrypt.cn", "账号")
	otpCmd.Flags().IntP("period", "p", 30, "刷新间隔")

	otpCmd.Flags().BoolP("qr", "", false, "生成二维码")
	otpCmd.Flags().IntP("width", "W", 200, "生成二维码宽度")
	otpCmd.Flags().IntP("height", "H", 200, "生成二维码高度")

	otpCmd.Flags().StringP("output", "o", "output.png", "输出图片")

	// otpCmd.Flags().BoolP("htop", "", false, "RFC4226, HOTP:An HMAC-Based One-Time Password Algorithm.")
	// otpCmd.Flags().BoolP("totp", "", false, "RFC6238, TOTP:Time-Based One-Time Password Algorithm.")

}

func NewOtpCmd() *cobra.Command {
	return otpCmd
}
