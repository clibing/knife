package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/h2non/filetype"
)

/**
 *JPG图片头信息:FFD8FF
 * PNG图片头信息:89504E47
 * GIF图片头信息:47494638
 * BMP图片头信息:424D
 */
func TestFileType(t *testing.T) {
	fmt.Println("测试File Type")

	fs, err := os.ReadDir("./img")
	if err != nil {
		fmt.Println("read dir error", err)
		return
	}

	for _, value := range fs {
		name := value.Name()
		buf, _ := os.ReadFile(fmt.Sprintf("./img/%s", name))

		kind, _ := filetype.Match(buf)
		if filetype.IsImage(buf) {
			fmt.Println("current file is image")
		}
		if kind == filetype.Unknown {
			fmt.Printf("Current img: img/%s, Unknown file type\n", name)
			continue
		}
		fmt.Printf("Current img: img/%s, file type: %s. MIME: %s\n", name, kind.Extension, kind.MIME.Value)
	}

}
