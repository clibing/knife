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
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"image/color"
)
var (
	size,level int
	fileName string
	disableBorder bool
)
// qrcodeCmd represents the qrcode command
var qrcodeCmd = &cobra.Command{
	Use:   "qrcode",
	Short: "生成二维码",
	Long: `将输入的内容生成二维码, 并生成png文件:

生成二维码
1. 当前目录快速生成二维码, 名字默认为 output.png
   knife qrcode "https://clibing.com"

2. 有边框，大小512，recovery level 2 输出到 /tmp/512.png 二维码的内容是 "https://clibing.com"
   knife qrcode -l 2 -s 512 -f /tmp/512.png "https://clibing.com"

3. 无边框，大小512，recovery level 2 输出到 /tmp/512.png 二维码的内容是 "https://clibing.com"
   knife qrcode -d -l 2 -s 512 -f /tmp/512.png "https://clibing.com"

.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Errorf("请输入二维码的内容")
		}
		recoveryLevel := qrcode.Medium
		if level == 2 {
			recoveryLevel = qrcode.High
		}else if level == 3 {
			recoveryLevel = qrcode.Highest
		}

		q, err := qrcode.New(args[0], recoveryLevel)
		if err != nil {
			fmt.Errorf("创建二维码基础参数异常 %s", err)
			return
		}
		q.BackgroundColor = color.White
		q.ForegroundColor = color.Black
		q.DisableBorder = disableBorder

		err = q.WriteFile(size, fileName)
		if err != nil {
			fmt.Errorf("生成二维码异常 %s", err)
		}
		fmt.Println("二维码生成成功, ", fileName)
	},
}

func init() {
	rootCmd.AddCommand(qrcodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// qrcodeCmd.PersistentFlags().String("foo", "", "A help for foo")
	qrcodeCmd.Flags().IntVarP(&size, "size", "s", 256, "二维码的size")
	qrcodeCmd.Flags().StringVarP(&fileName, "fileName", "f", "./output.png", "输出二维码图片完整路径")

	qrcodeCmd.Flags().IntVarP(&level, "level", "l", 1, "生成图片质量，默认是1(15%), 2(25%), 3(30%)")
	qrcodeCmd.Flags().BoolVarP(&disableBorder, "disableBorder", "d", false, "是否禁用二维码边框, 默认不禁用")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// qrcodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
