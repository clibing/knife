package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
)

// base64Cmd represents the image command
var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "图片转base64",
	Long: `将图片转换成base64， 也可以将base64生成图片:

1. 图片转base64
  knife image base64 -i  ./1.png -o ./result.txt

2. base64转图片
  knife image -b=false -i ./1.txt -o ./1.png

3. 图片转base64后，生成html，用于预览
  --html 生成html
  -o 将结果输出到指定文件
  1.png 2.png 3.webp 输入的内容
  knife image base64 --html -o out.txt 1.png 2.jpg 3.webp`,
	Run: func(c *cobra.Command, args []string) {
		input, _ := c.Flags().GetString("input")
		output, _ := c.Flags().GetString("output")
		binary, _ := c.Flags().GetBool("binary")

		if len(args) == 0 && len(input) == 0 {
			c.Help()
			return
		}

		if binary {
			toBase64(c, args, input, output)
		} else {
			toImage(c, args, input, output)
		}
	},
}

/**
 * base64 --> image
 */
func toImage(c *cobra.Command, args []string, input string, output string) {
	if len(args) > 0 {
		for i, v := range args {
			o := fmt.Sprintf("%d_%s", i, output)
			file, err := os.ReadFile(v)
			if err != nil {
				fmt.Println("load base64 failed: ", err.Error())
				continue
			}
			result, err := base64.StdEncoding.DecodeString(string(file))
			if err != nil {
				fmt.Println("to base64 failed: ", err.Error())
				continue
			}
			os.WriteFile(o, result, 0664)
		}
	} else if len(input) > 0 {
		file, err := os.ReadFile(input)
		if err != nil {
			fmt.Println("load base64 failed: ", err.Error())
			return
		}
		result, err := base64.StdEncoding.DecodeString(string(file))
		if err != nil {
			fmt.Println("decode failed: ", err.Error())
			return
		}
		os.WriteFile(output, []byte(result), 0664)
	} else {
		fmt.Println("输入为空")
	}
}

/**
 * 获取文件的 mime type
 */
func getMimeType(image, result []byte) string {
	m := mimetype.Detect(image)
	return fmt.Sprintf("data:%s;base64,%s", m, string(result))
}

/**
 * image --> base64
 */
func toBase64(c *cobra.Command, args []string, input string, output string) {

	data := make([]string, 0)

	if len(args) > 0 {
		for i, v := range args {
			o := fmt.Sprintf("%d_%s", i, output)
			file, err := os.ReadFile(v)
			if err != nil {
				fmt.Println("load image failed: ", err.Error())
				continue
			}
			result := base64.StdEncoding.EncodeToString(file)
			os.WriteFile(o, []byte(result), 0664)
			v := getMimeType(file, []byte(result))
			data = append(data, v)
		}
	} else if len(input) > 0 {
		file, err := os.ReadFile(input)
		if err != nil {
			fmt.Println("load image failed: ", err.Error())
			return
		}
		result := base64.StdEncoding.EncodeToString(file)
		os.WriteFile(output, []byte(result), 0664)
		v := getMimeType(file, []byte(result))
		data = append(data, v)
	} else {
		fmt.Println("输入为空")
	}
	if len(data) > 0 {
		html(c, fmt.Sprintf("%s.html", output), data...)
	}
}

func html(c *cobra.Command, output string, args ...string) {
	need, _ := c.Flags().GetBool("html")
	if !need {
		return
	}

	template := template.New(output)
	t, err := template.Parse(`
<html>
<body>

{{ range .data }}
	<img src="{{.}}"/>
{{ end }}

</body>
</html>	

	`)
	if err != nil {
		fmt.Println("parse html failed: ", err.Error())
		return
	}
	var html bytes.Buffer
	data := make(map[string]interface{}, 1)
	data["data"] = args
	t.Execute(&html, data)
	os.WriteFile(output, html.Bytes(), 0664)
}

func init() {
	base64Cmd.Flags().StringP("input", "i", "", "源文件")
	base64Cmd.Flags().StringP("output", "o", "", "输出到文件")
	base64Cmd.Flags().BoolP("html", "", false, "图片转base64后，是否生成html，默认不需要")
	base64Cmd.Flags().BoolP("binary", "b", true, "是否为二进制, 默认图片转base64")
}

func NewImageCmd() *cobra.Command {
	return base64Cmd
}
