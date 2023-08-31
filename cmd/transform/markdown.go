package transform

import (
	"fmt"

	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/spf13/cobra"
)

// mdCmd represents the md command
var mdCmd = &cobra.Command{
	Use:     "markdown",
	Aliases: []string{"md"},
	Short:   "markdown html互转工具",
	Example: `1. source.html内容: <html><body>hello 'clibing' </br><a href='https://github.com/clibing'> clibing </a></body></html>
2. html -> markdown
源文件到目标文件
knife convert markdown -s source.html -t target.md

输出到控制台
knife convert markdown "<html><body><h1>Hello World</h1><img src=\"/static/logo.png\" /></body></html>"

输出到文件
knife convert markdown "<html><body><h1>Hello World</h1><img src=\"/static/logo.png\" /></body></html>" -t ./out.md

3. source.md 内容
# hello 'clibing' 

[clibing](https://github.com/clibing)

4. markdown to html
knife convert markdown -d html -s source.md 
knife convert markdown -d html "# hello 'clibing'"
knife convert markdown -d html -s source.md -t target.html`,
	Run: func(c *cobra.Command, args []string) {
		data := make([]string, 0)
		source, e := c.Flags().GetString("source")
		var content string
		if e == nil && len(source) > 0 {
			value, e := os.ReadFile(source)
			if e != nil {
				fmt.Println("读取源文件异常: ", e.Error())
			}
			content = string(value)
			data = append(data, string(content))
		}

		if len(content) == 0 && len(args) == 0 {
			fmt.Println("暂无输入源")
			return
		}

		if len(args) > 0 {
			data = append(data, args...)
		}

		// direct 默认为false
		direct, _ := c.Flags().GetString("direct")
		for _, d := range data {
			var value []byte
			if direct == "html" {
				value = toHtml([]byte(d))
			} else {
				domain, _ := c.Flags().GetString("domain")
				value = toMd(domain, []byte(d))
			}
			target, e := c.Flags().GetString("target")
			if e == nil && len(target) > 0 {
				os.WriteFile(target, value, 0644)
			}
			if len(target) == 0 {
				fmt.Println(string(value))
			}
		}
	},
}

func toMd(domain string, data []byte) []byte {
	converter := md.NewConverter(domain, true, nil)
	markdown, err := converter.ConvertBytes(data)
	if err != nil {
		return nil
	}
	return markdown
}

/**
 * convert to html
 */
func toHtml(data []byte) []byte {
	htmlFlags := html.CommonFlags | html.CompletePage | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.ToHTML(data, nil, renderer)
}

func init() {
	mdCmd.Flags().StringP("source", "s", "", "源文件")
	mdCmd.Flags().StringP("target", "t", "", "目标文件")
	mdCmd.Flags().StringP("domain", "D", "", "html转markdown，用于拼接html中img: src=\"/image/logo.png\" --> ![](${domain}/image/logo.png)")
	mdCmd.Flags().StringP("direct", "d", "", "转换的目标类型: html|md; 默认: html到markdown")
}

func NewMarkdown() *cobra.Command {
	return mdCmd
}
