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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	action string
	direct bool
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "常用的加密算法",
	Long: `封装一些常用的加密算法包括md5, sha, sha256, hmac, ras, aes等等: 

1. MD5 使用方式
   knife sign -md5 ""
   knife sign -md5 -file /tmp/file.txt

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, content := range args {
			switch action {
			// md5
			case "md5":
				{
					h := md5.New()
					h.Write([]byte(content))
					value := hex.EncodeToString(h.Sum(nil))
					fmt.Println("source: ", content)
					fmt.Println("md5   : ", value)
				}
				break
			case "sha1":
				{
					s := sha1.New()
					s.Write([]byte(content))
					value := hex.EncodeToString(s.Sum(nil))
					fmt.Println("source: ", content)
					fmt.Println("sha1  : ", value)
				}
				break
			case "sha256":
				{
					s := sha256.New()
					s.Write([]byte(content))
					value := hex.EncodeToString(s.Sum(nil))
					fmt.Println("source: ", content)
					fmt.Println("sha256: ", value)
				}
				break
			case "sha512":
				{
					s := sha512.New()
					s.Write([]byte(content))
					value := hex.EncodeToString(s.Sum(nil))
					fmt.Println("source: ", content)
					fmt.Println("sha512: ", value)
				}
				break
			case "base64":
				{
					// 加密
					if direct == false {
						v := base64.StdEncoding.EncodeToString([]byte(content))
						fmt.Println("source: ", content)
						fmt.Println("base64: ", v)
					} else {
						v, _ := base64.StdEncoding.DecodeString(content)
						fmt.Println("base64: ", content)
						fmt.Println("source: ", string(v))
					}
				}
				break
			default:
				fmt.Println("暂不支持的加密方法, ", action)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")
	signCmd.Flags().StringVarP(&action, "type", "t", "", "选择加密方法")
	signCmd.Flags().BoolVarP(&direct, "direct", "d", false, "加密或者解密,支持部分可逆的算法, 默认加密")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
