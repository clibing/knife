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
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	inputTimeNum, inputTimeStr bool   // 输入的内容是否为数字类型、字符串类型（默认为当前时间）, 输入是否为格式化的字符串(默认为数字)
	duration, format           string // duration 需要累加时间; format格式化
)
/**
* 1. 先获取时间
*    a. 默认为 time.Now()
     b. -n 输入的为 long
     c. -s 输入的为 字符串
* 2. -d 是否偏移  such as "300ms", "-1.5h" or "2h45m".
* 3. 输出
*    a. 输出的格式 数字还是字符串
*    a. 格式化样式
*    b. 输出样式
*/

// timeCmd represents the unixTime command
var timeCmd = &cobra.Command{
	Use: "time",
	//TraverseChildren: true,
	Short: "时间小工具",
	Long: `在开发中经常使用时间转换判断大小:

时间格式化样式:
	"yyyyMMddHHmmss"
	"yyyy-MM-dd HH:mm:ss"
	"yyyy.MM.dd HH:mm:ss"
	"yyyy/MM/dd HH:mm:ss"
	"yyyyMMdd HH:mm:ss"
	"yyyyMMdd HH:mm:ss"
	"..."

具体使用：

1. 快速获取当前的时间戳13位、10位和当前默认格式化的时间
knife time
输出: 
  毫秒: 1647315266451
  秒　: 1647315266
  美化: 2022-03-15 11:34:26

2. 接收的毫秒、秒进行格式化
knife time -n 1636284416438 
输出: 
  毫秒: 1636284416438
  秒　: 1636284416
  美化: 2021-11-07 19:26:56

knife time -n 1636284416000 
输出: 
  毫秒: 1636284416000
  秒　: 1636284416
  美化: 2021-11-07 19:26:56

2. 接收格式化后的时间进行转化，默认为yyyy-MM-dd HH:mm:ss格式，如果输入的不是这个格式，需要手动指定 
knife time -s "2022-03-16 11:23:10"  
输出: 
  毫秒: 1647400990000
  秒　: 1647400990
  美化: 2022-03-16 11:23:10

knife time -s "2022/03/16 11:23:10" -f "yyyy/MM/dd HH:mm:ss"
输出: 
  毫秒: 1647400990000
  秒　: 1647400990
  美化: 2022/03/16 11:23:10

3. 偏移当前时间或者指定的时间
knife time -d 1h  偏移当前时间1个小时以后，也可以是 1h30m, 1h10s                                       
输出: 
  毫秒: 1647319176467
  秒　: 1647319176
  美化: 2022-03-15 12:39:36

4. 复杂的使用
knife time -s "2022/03/16 11:23:10" -f "yyyy/MM/dd HH:mm:ss" -d 1h30s
输出: 
  毫秒: 1647404620000
  秒　: 1647404620
  美化: 2022/03/16 12:23:40

knife  time -s -f "yyyy/MM/dd HH:mm:ss" -d 1h30s  "2022/03/16 11:23:10"
输出: 
  毫秒: 1647404620000
  秒　: 1647404620
  美化: 2022/03/16 12:23:40

2022/03/16 11:23:10
2022/03/16 12:23:40

5. 注意 不要同事接收数字类型和字符串类型
`,
	Run: func(_ *cobra.Command, args []string) {
		var value string
		if len(args) >= 1 {
			value = args[0]
		}
		if inputTimeNum && inputTimeStr {
			fmt.Println("不能同时接收时间戳和格式化后的时间")
			return
		}

		t := time.Now()
		var v int64
		// 输入的为数字类型
		if inputTimeNum {
			len := len(value)
			var e error
			v, e = strconv.ParseInt(value, 10, 64)
			if e != nil {
				fmt.Printf("输入[%s]的数字类型的时间戳异常, msg: %s\n", value, e.Error())
				return
			}
			if len == 13 {
				t = time.UnixMilli(v)
			}else if len == 10{
				t = time.Unix(v, 0)
			}else{
				fmt.Printf("输入[%d]的时间戳不是13位或者10位, len: %d\n", v, len)
				return
			}
		}

		// 生成格式化 样式
		f := convertFormat(format)
		if inputTimeStr {
			var e error
			t, e = time.ParseInLocation(f, value, time.Local)
			if e != nil {
				fmt.Printf("输入[%s]的字符串时间，格式化[%s], msg: %s\n", value, f, e.Error())
				return
			}
		}

		if len(duration) > 0 {
			pd, err := time.ParseDuration(duration)
			if err != nil {
				fmt.Printf("偏移时间异常(参考: 1h10m, 1h10s, 1.5h30m等): %s, msg: %s\n", duration, err.Error())
				return
			}
			t = t.Add(pd)
		}

		fmt.Printf("输出: \n  毫秒: %d\n  秒　: %d\n  美化: %s\n", t.UnixMilli(), t.Unix(), t.Format(f))
	},
}

func convertFormat(format string) string {
	var f string
	if len(format) != 0 {
		f = strings.Replace(format, "yyyy", "2006", 1)
		f = strings.Replace(f, "MM", "01", 1)
		f = strings.Replace(f, "dd", "02", 1)
		f = strings.Replace(f, "HH", "15", 1)
		f = strings.Replace(f, "mm", "04", 1)
		f = strings.Replace(f, "ss", "05", 1)
	} else {
		f = "2006-01-02 15:04:05"
	}
	return f
}

func init() {
	rootCmd.AddCommand(timeCmd)
	timeCmd.Flags().BoolVarP(&inputTimeNum, "inputTimeNum", "n", false, "开启，输入的内容为数字类型时间戳，自动识别13、10位")
	timeCmd.Flags().BoolVarP(&inputTimeStr, "inputTimeStr", "s", false, "开启，输入的内容为格式化的时间戳, 默认 yyyy-MM-dd HH:mm:ss")
	timeCmd.Flags().StringVarP(&duration, "duration", "d", "", "时间偏转计算(可以增加小时、分钟、秒、毫秒)，例如\"300ms\", \"-1.5h\" or \"2h45m\"")
	timeCmd.Flags().StringVarP(&format, "format", "f", "", "时间格式化, 默认为 yyyy-MM-dd HH:mm:ss")
}
