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
	Long: `对URL进行编码、解码：

编码：
knife url  "http://github.com/clibing/knife"
http%3A%2F%2Fgithub.com%2Fclibing%2Fknife

解码：
knife url -e "http%3A%2F%2Fgithub.com%2Fclibing%2Fknife"
http://github.com/clibing/knife

支持管道
编码： echo "http://github.com/clibing/knife" | knife url 
解码： echo "http%3A%2F%2Fgithub.com%2Fclibing%2Fknife" | knife url -e
`,
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
	rootCmd.AddCommand(urlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// urlCmd.PersistentFlags().String("foo", "", "A help for foo")
	urlCmd.Flags().BoolVarP(&unescape, "unescape", "e", false, "URL编码, 默认编码")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// urlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
