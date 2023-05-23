/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	htmlToMd "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)


var (
	source, target string
)
// mdCmd represents the md command
var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "markdown html互转工具",
	Long: `支持markdown->html, 或者html->markdown:

1. source.html内容
<html>
	<body>
		hello 'clibing' </br>
		<a href='https://github.com/clibing'> clibing </a>
	</body>
</html>
2. html -> markdown
knife md -s source.html -t target.md

3. source.md 内容
hello 'clibing' 

[clibing](https://github.com/clibing)

4. markdown to html
knife md -d -s source.md -t target.html.`,
	Run: func(_ *cobra.Command, _ []string) {
		// direct 默认为false
		if direct {
			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)
			md, err := ioutil.ReadFile(source)
			if err != nil {
				log.Println(err.Error())
			}
			html := markdown.ToHTML(md, nil, renderer)
			ioutil.WriteFile(target, html, 0644)
		}else{
			converter := htmlToMd.NewConverter("", true, nil)
			html, err := ioutil.ReadFile(source)
			if err != nil {
				log.Fatal(err)
				return
			}
			markdown, err := converter.ConvertString(string(html))
			if err != nil {
				log.Fatal(err)
				return
			}
			err = ioutil.WriteFile(target, []byte(markdown), 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mdCmd.PersistentFlags().String("foo", "", "A help for foo")
	mdCmd.Flags().StringVarP(&source, "source", "s", "", "源文件")
	mdCmd.Flags().StringVarP(&target, "target", "t", "", "目标文件")
	mdCmd.Flags().BoolVarP(&direct, "direct", "d", false, "转换的目标类型，from->to，默认是HTML到markdown")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
