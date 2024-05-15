package image

import (
	"fmt"
	"testing"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func TestImage(t *testing.T) {
	imagick.Initialize()
	defer imagick.Terminate()
	var err error
	// 加载底图
	logo := imagick.NewMagickWand()
	err = logo.ReadImage("input.png")
	if err != nil {
		fmt.Println("读取文件Logo发生错错误", err)
		return
	}

	w := logo.GetImageWidth()
	h := logo.GetImageHeight()

	logo.ResizeImage(w/2, h/2, imagick.FILTER_LANCZOS)
	logo.SetImageCompressionQuality(80)

	logo.WriteImage("output.jpg")

	logo.Destroy()
	imagick.Terminate()
}
