/* Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	f,h,d,v string
	z int8
	i uint64
	s bool
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
注意：格式化有点另类，不是java的yyyyMMddHHmmss这样的格式，是按照上面的数字样式格式化的，注意数值不能变

其中年月日在中间可以增加分隔符，顺序可以替换， 时间部分同理， 其中2006、05可以省略
各个参数的级别:
	-i(--input) > -s(--second)

具体使用：

1. 获取13为的时间戳
knife time
1636367686555

2. 获取10为秒级的时间戳
knife time -s
1636367686

3. 按照格式格式化
knife time -f "2006-01-02 15:04:05"
2021-11-08 18:36:08

4. 获取当前时间之后的1小时，并按照格式样式格式化
knife time -f "2006-01-02 15:04:05" -d 1h
2021-11-08 19:38:28

5. 时间戳格式化
knife time -i 1636367761724   
2021-11-08 18:36:01 # 默认是  "2006-01-02 15:04:05"

6. 指定样式格式化时间戳
knife time -i 1636367761724 -f "2006-01-02 15:04"  
2021-11-08 18:41
`,
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		if d != "" {
			pd ,_ := time.ParseDuration(d)
			t = t.Add(pd)
		}
		if f != "" {
			fmt.Println(t.Format(f))
			return
		}

		if i > 0 {
			if f == "" {
				convertTime(i, "2006-01-02 15:04:05")
			} else{
				convertTime(i, f)
			}
			return
		}

		if v != "" {
			stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
			if s == false {
				fmt.Println(stamp.UnixMilli())
				return
			}
			fmt.Println(stamp.Unix())
			return
		}

		if s == false {
			fmt.Println(t.UnixMilli())
			return
		} else {
			fmt.Println(t.Unix())
			return
		}
	},
}

func convertTime(uTimeMills uint64, format string) string {
	v := time.Unix(int64(uTimeMills/1000), 0).Format(format)
	fmt.Println(v)
	return v
}

func init() {
	rootCmd.AddCommand(timeCmd)

	// Here you will define your flags and configuration settings.
	timeCmd.Flags().StringVarP(&f, "format", "f", "", "时间格式化")
	timeCmd.Flags().StringVarP(&d, "duration", "d", "", "时间的偏移, 支持 (1h, -1h, 1m, -1m)")
	//timeCmd.Flags().Int8VarP(&z, "zone", "z", 8, "时区")
	timeCmd.Flags().Uint64VarP(&i, "input", "i", 0, "需要格式化的int类型的时间戳")
	timeCmd.Flags().BoolVarP(&s, "second", "s", false, "是否需要秒, 默认不需要, 暂时毫秒")
	timeCmd.Flags().StringVarP(&v, "value", "v", "", "输入的字符串时间格式为: yyyy-MM-dd HH:mm:ss")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
