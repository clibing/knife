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
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
	"os"
	"strings"
)

var (
	prefix, indent string
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "json美化器, 支持管道输入",
	Long: `在传输过程中一般会对json压缩成一行，但是在查看的时候非常不友好, 
每次都需要打开在线格式化工具，用着不爽，故出现此工具， 默认采用 '\t' 进行缩进:

例如：
原json信息 

{"id":1,"name":"clibing"}

输出新的json信息
{
    "id": 1,
    "name": "clibing"
}

echo '{"id":1,"name":"clibing"}' | go run main.go json

如果此格式化不能满足 可以安装jq进行格式化
jq: https://github.com/stedolan/jq.git
使用:
echo '{"id":1,"name":"clibing"}' | jq.`,
	Run: func(cmd *cobra.Command, args []string) {
		option := &pretty.Options{Width: 80, Prefix: prefix, Indent: indent, SortKeys: false}
		for _, json := range args {
			fmt.Printf("%s\n", pretty.PrettyOptions([]byte(json), option))
		}
		if len(args) > 0 {
			return
		}

		fileInfo, _ := os.Stdin.Stat()
		if (fileInfo.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
			fmt.Println("The command is intended to work with pipes.")
			return
		}

		var buf strings.Builder
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			buf.WriteString(s.Text())
		}
		fmt.Printf("%s\n", pretty.PrettyOptions([]byte(buf.String()), option))
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)

	// Here you will define your flags and configuration settings.
	jsonCmd.Flags().StringVarP(&indent, "indent", "i", "\t", "json格式化缩进标记，默认制表符")
	jsonCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "json格式化前缀")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
