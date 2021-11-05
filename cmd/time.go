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
	"time"

	"github.com/spf13/cobra"
)

var (
	f,h string
	z,i int
	m bool
)
// timeCmd represents the unixTime command
var timeCmd = &cobra.Command{
	Use:   "time",
	//TraverseChildren: true,
	Short: "时间小工具",
	Long: `在开发中经常使用时间转换判断大小:
时间格式化样式:
	"20060102150405"
	"200601021504"
	"01021504"
	"0102150405"
	"2006-01-02"
	"2006-01-02 15:04"
	"2006-01-02 15:04:05"
	"2006/01/02 15:04:05"
	"01/02/2006 15:04:05"
	"20060102 15:04:05"
其中年月日在中间可以增加分隔符，顺序可以替换， 时间部分同理， 其中2006、05可以省略`,
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		if f != "" {
			fmt.Println(t.Format(f))
			return
		}
		if m == true {
			fmt.Println(t.UnixMilli())
			return
		} else {
			fmt.Println(t.Unix())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

	// Here you will define your flags and configuration settings.
	timeCmd.Flags().StringVarP(&f, "format", "f", "", "时间格式化")
	timeCmd.Flags().IntVarP(&z, "zone", "z", 8, "时区")
	timeCmd.Flags().IntVarP(&i, "input", "i", 0, "需要格式化的int类型的时间戳")
	timeCmd.Flags().BoolVarP(&m, "milli", "m", true, "是否需要毫秒, 默认需要")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
