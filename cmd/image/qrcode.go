package image

import (
	"fmt"
	"image/png"
	"os"

	// 	"github.com/skip2/go-qrcode"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/tuotoo/qrcode"

	"github.com/spf13/cobra"
)

// qrcodeCmd represents the qrcode command
var qrcodeCmd = &cobra.Command{
	Use:   "qrcode",
	Short: "生成二维码",
	Long: `将输入的内容生成二维码, 并生成png文件:

生成二维码
1. 当前目录快速生成二维码, 名字默认为 output.png
   knife qrcode "https://clibing.com"

2. 有边框，大小512，recovery level 2 输出到 /tmp/512.png 二维码的内容是 "https://clibing.com"
   knife qrcode -l 2 -s 512 -f /tmp/512.png "https://clibing.com"

3. 无边框，大小512，recovery level 2 输出到 /tmp/512.png 二维码的内容是 "https://clibing.com"
   knife qrcode -d -l 2 -s 512 -f /tmp/512.png "https://clibing.com"

.`,
	Run: func(c *cobra.Command, args []string) {
		decode, _ := c.Flags().GetBool("decode")
		if decode {
			input, _ := c.Flags().GetString("input")
			fmt.Println("input:", input)
			fi, err := os.Open(input)
			if err != nil {
				fmt.Println("打开文件异常", err.Error())
				return
			}
			defer fi.Close()
			matrix, err := qrcode.Decode(fi)
			if err != nil {
				fmt.Println("解析二维码内容失败", err.Error())
				return
			}
			fmt.Println(matrix.Content)
			return
		}

		output, _ := c.Flags().GetString("output")
		width, _ := c.Flags().GetInt("width")
		height, _ := c.Flags().GetInt("height")

		// Create the barcode
		qrCode, _ := qr.Encode(args[0], qr.M, qr.Auto)

		// Scale the barcode to 200x200 pixels
		qrCode, _ = barcode.Scale(qrCode, width, height)

		// create the output file
		file, _ := os.Create(output)
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, qrCode)

		fmt.Println("二维码生成成功, ", output)
	},
}

func init() {
	qrcodeCmd.Flags().StringP("key", "k", "", "二维码内容")
	qrcodeCmd.Flags().IntP("width", "W", 256, "二维码宽度")
	qrcodeCmd.Flags().IntP("height", "H", 256, "二维码高度")
	qrcodeCmd.Flags().StringP("output", "o", "./output.png", "输出二维码图片完整路径")

	qrcodeCmd.Flags().BoolP("decode", "d", false, "识别二维码内容")
	qrcodeCmd.Flags().StringP("input", "i", "", "二维码图片源地址")

}

func NewQrcodeCmd() *cobra.Command {
	return qrcodeCmd
}
