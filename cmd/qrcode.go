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
	size int
	name string
)
// qrcodeCmd represents the qrcode command
var qrcodeCmd = &cobra.Command{
	Use:   "qrcode",
	Short: "生成二维码",
	Long: `将输入的内容生成二维码, 并生成png文件:

生成二维码的内容为"admin"， 名字为 image
knife qrcode -s 128 -n image "admin"
.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Errorf("请输入二维码的内容")
		}
		err := qrcode.WriteColorFile(args[0], qrcode.Medium, size, color.Black, color.White, name+".png")
		if err != nil {
			fmt.Errorf("create qrcode error: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(qrcodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// qrcodeCmd.PersistentFlags().String("foo", "", "A help for foo")
	qrcodeCmd.Flags().IntVarP(&size, "size", "s", 256, "二维码的size")
	qrcodeCmd.Flags().StringVarP(&name, "name", "n", "qr.png", "二维码图片的名字")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// qrcodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
