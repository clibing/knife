package image

import (
	"fmt"
	"image/color"
	"log"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

var (
	size, level   int
	fileName      string
	disableBorder bool
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
	Run: func(_ *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("请输入二维码的内容")
			return
		}
		recoveryLevel := qrcode.Medium
		if level == 2 {
			recoveryLevel = qrcode.High
		} else if level == 3 {
			recoveryLevel = qrcode.Highest
		}

		q, err := qrcode.New(args[0], recoveryLevel)
		if err != nil {
			log.Fatalf("创建二维码基础参数异常 %s", err)
			return
		}
		q.BackgroundColor = color.White
		q.ForegroundColor = color.Black
		q.DisableBorder = disableBorder

		err = q.WriteFile(size, fileName)
		if err != nil {
			log.Fatalf("生成二维码异常 %s", err)
			return
		}
		fmt.Println("二维码生成成功, ", fileName)
	},
}

func init() {
	qrcodeCmd.Flags().IntVarP(&size, "size", "s", 256, "二维码的size")
	qrcodeCmd.Flags().StringVarP(&fileName, "fileName", "f", "./output.png", "输出二维码图片完整路径")

	qrcodeCmd.Flags().IntVarP(&level, "level", "l", 1, "生成图片质量，默认是1(15%), 2(25%), 3(30%)")
	qrcodeCmd.Flags().BoolVarP(&disableBorder, "disableBorder", "d", false, "是否禁用二维码边框, 默认不禁用")
}

func NewQrcodeCmd() *cobra.Command {
	return qrcodeCmd
}
