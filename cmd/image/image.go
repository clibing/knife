package image

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	imageFile, base64File string
	direct                bool
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "图片转base64",
	Long: `将图片转换成base64， 也可以将base64生成图片:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, _ []string) {
		if direct {
			// base64 -> image
			ddd, _ := base64.StdEncoding.DecodeString(base64File) //成图片文件并把文件写入到buffer
			err := os.WriteFile(imageFile, ddd, 0666)             //buffer输出到jpg文件中（不做处理，直接写到文件）
			if err != nil {
				fmt.Printf("base64 to image error, %s", err)
				return
			}
		} else {
			// image -> base64
			file, err := os.ReadFile(imageFile)
			if err != nil {
				fmt.Printf("read file error, %s", err)
				return
			}
			result := base64.StdEncoding.EncodeToString(file)
			os.WriteFile(base64File, []byte(result), 0664)
		}
	},
}

func init() {
	imageCmd.Flags().StringVarP(&imageFile, "imageFile", "i", "", "图片文件")
	imageCmd.Flags().StringVarP(&base64File, "base64File", "b", "", "base64存储的文件")
	imageCmd.Flags().BoolVarP(&direct, "direct", "d", false, "base64生成图片，默认false， 图片生成base64")
}

func NewImageCmd() *cobra.Command {
	return imageCmd
}
