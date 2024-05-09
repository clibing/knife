package qrcode

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func TestQrCode(t *testing.T) {
	// open and decode image file
	file, _ := os.Open("/Users/liubaixun/Downloads/macOS-collect/otp.png")
	img, s, _ := image.Decode(file)
	fmt.Println(s)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	fmt.Println(result)
}
