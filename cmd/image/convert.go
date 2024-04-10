package image

/**
 * 图片转换工具类
 */

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/spf13/cobra"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// convertCmd represents the qrcode command
var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"c"},
	Short:   "图片格式转换",
	Long: `图片转换: 

1. 将svg转换为png，并指定生成的图片100x100
    knife image convert -t 0 -i 1.svg -o 1.png -w 100 -h 100
.`,
	Run: func(cmd *cobra.Command, args []string) {
		_type, _ := cmd.Flags().GetInt("type")
		w, _ := cmd.Flags().GetInt("width")
		h, _ := cmd.Flags().GetInt("height")
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")

		switch _type {
		case 0:
			e := SvgToPng(input, output, w, h)
			if e != nil {
				fmt.Printf("转换失败: %s\n", e.Error())
			}
		}
	},
}

func init() {
	convertCmd.Flags().StringP("input", "i", "", "输入资源文件路径")
	convertCmd.Flags().StringP("output", "o", "", "输出资源文件路径")
	convertCmd.Flags().IntP("width", "W", 0, "输出图片宽度")
	convertCmd.Flags().IntP("height", "H", 0, "输出图片高度")
	convertCmd.Flags().IntP("type", "t", 0, `转换类型:
    0: svg->png, 默认输入input.svg 输出output.png;
    1: png->jpg, 默认输入input.png 输出output.jpg;
    2: jpg->png, 默认输入input.jpg 输出output.png;
	`)
}

func NewConvertCmd() *cobra.Command {
	return convertCmd
}

func SvgToPng(input, output string, w, h int) error {
	if len(input) == 0 {
		input = "input.svg"
	}

	if len(output) == 0 {
		output = "output.png"
	}

	in, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("打开svg文件错误: %s", err.Error())
	}
	defer in.Close()

	icon, _ := oksvg.ReadIconStream(in)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("生成输出文件错误: %s", err.Error())
	}
	defer out.Close()

	err = png.Encode(out, rgba)
	if err != nil {
		return fmt.Errorf("转换异常: %s", err.Error())
	}
	return nil
}
