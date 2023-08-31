package transform

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// 是否
	unescape bool
)

// urlCmd represents the url command
var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "URL的编码 解码",
	Example: `对URL进行编码、解码：

编码：
knife convert url "http://github.com/clibing/knife"
http%3A%2F%2Fgithub.com%2Fclibing%2Fknife

解码：
knife convert url -e "http%3A%2F%2Fgithub.com%2Fclibing%2Fknife"
http://github.com/clibing/knife

支持管道
编码： echo "http://github.com/clibing/knife" | knife convert url 
解码： echo "http%3A%2F%2Fgithub.com%2Fclibing%2Fknife" | knife convert url -e`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, value := range args {
				if !unescape {
					fmt.Printf("%s\n", url.QueryEscape(value))
				} else {
					result, _ := url.QueryUnescape(value)
					fmt.Printf("%s\n", result)
				}
			}
		}
		value, _ := os.Stdin.Stat()
		if (value.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
			return
		}

		var buf strings.Builder
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			buf.WriteString(s.Text())
		}
		if !unescape {
			fmt.Printf("%s\n", url.QueryEscape(buf.String()))
		} else {
			result, _ := url.QueryUnescape(buf.String())
			fmt.Printf("%s\n", result)
		}
	},
}

func init() {
	urlCmd.Flags().BoolVarP(&unescape, "unescape", "e", false, "URL编码, 默认编码")
}

func NewUrlEncoding() *cobra.Command {
	return urlCmd
}
