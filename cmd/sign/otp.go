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
	Run: func(c *cobra.Command, args []string) {
		qr, _ := c.Flags().GetBool("qr")
		account, _ := c.Flags().GetString("account")
		issuer, _ := c.Flags().GetString("issuer")
		secret, _ := c.Flags().GetString("secret")

		// 生成秘钥
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      issuer,
			AccountName: account,
			Secret:      []byte(secret),
		})

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
			height, _ := c.Flags().GetInt("output")

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
