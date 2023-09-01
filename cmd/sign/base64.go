package sign

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: `base64 签名计算`,
	Run: func(c *cobra.Command, args []string) {
		input, _ := c.Flags().GetString("input")
		if len(args) == 0 && len(input) == 0 {
			c.Help()
		}
		if len(input) > 0 {
			value, err := os.ReadFile(input)
			if err != nil {
				fmt.Println("sign file error, ", err)
			} else {
				base64method(c, value)
			}
		}
		if len(args) > 0 {
			for _, content := range args {
				base64method(c, []byte(content))
			}
		}
	},
}

func base64method(c *cobra.Command, content []byte) {
	direct, _ := c.Flags().GetBool("direct")
	output, _ := c.Flags().GetString("output")

	skip_write := len(output) == 0

	// 加密
	var result, value string
	if !direct {
		if skip_write {
			fmt.Println("source :", strings.Replace(string(content), "\n", "", -1))
		}
		value = base64.StdEncoding.EncodeToString(content)
		result = "base64"
	} else {
		if skip_write {
			fmt.Println("base64 :", strings.Replace(string(content), "\n", "", -1))
		}
		data, e := base64.StdEncoding.DecodeString(string(content))
		if e != nil {
			fmt.Println("解码异常: ", e.Error())
			return
		}
		value = string(data)
		result = "source"
	}

	if skip_write {
		fmt.Printf("%s: %s\n", result, value)
	} else {
		os.WriteFile(output, []byte(value), 0644)
	}
}

func init() {
	base64Cmd.Flags().StringP("input", "i", "", "输入文件")
	base64Cmd.Flags().StringP("output", "o", "", "输出的文件")
	base64Cmd.Flags().BoolP("direct", "d", false, "编码方式，默认encode")
}

func NewBase64Cmd() *cobra.Command {
	return base64Cmd
}
