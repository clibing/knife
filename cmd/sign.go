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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var (
	action, sourceFile string
	direct             bool
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "常用的加密算法",
	Long: `封装一些常用的加密算法包括md5, sha, sha256, hmac, ras, aes等等: 

使用方式
knife sign -t <type> "" 其中: <type>支持"md5", "sha1", "sha256", "sha512", "base64"

1. MD5 
   knife sign -t md5 "clibing"
   knife sign -t md5 -s /tmp/data.txt 注意文件签名与指定字符串签名不一致， 因为文件最后含有一个\r\n 、\r之类的换行符是隐藏的
   echo "clibing" | knife sign -t md5
2. sha1, sha256, sha512, base64操作同md5.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, content := range args {
			sign(action, []byte(content))
		}
		if sourceFile != "" {
			value, err := ioutil.ReadFile(sourceFile)
			if err != nil {
				fmt.Println("sign file error, ", err)
			}
			fmt.Println("source file: ", sourceFile)
			sign(action, value)
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
		sign(action, []byte(buf.String()))
	},
}

func sign(signType string, content []byte) {
	switch signType {
	// md5
	case "md5":
		{
			h := md5.New()
			h.Write(content)
			value := hex.EncodeToString(h.Sum(nil))
			fmt.Println("source: ", string(content))
			fmt.Println("md5   : ", value)
		}
		break
	case "sha1":
		{
			s := sha1.New()
			s.Write(content)
			value := hex.EncodeToString(s.Sum(nil))
			fmt.Println("source: ", string(content))
			fmt.Println("sha1  : ", value)
		}
		break
	case "sha256":
		{
			s := sha256.New()
			s.Write(content)
			value := hex.EncodeToString(s.Sum(nil))
			fmt.Println("source: ", string(content))
			fmt.Println("sha256: ", value)
		}
		break
	case "sha512":
		{
			s := sha512.New()
			s.Write(content)
			value := hex.EncodeToString(s.Sum(nil))
			fmt.Println("source: ", string(content))
			fmt.Println("sha512: ", value)
		}
		break
	case "base64":
		{
			// 加密
			if direct == false {
				v := base64.StdEncoding.EncodeToString(content)
				fmt.Println("source: ", string(content))
				fmt.Println("base64: ", v)
			} else {
				v, _ := base64.StdEncoding.DecodeString(string(content))
				fmt.Println("base64: ", string(content))
				fmt.Println("source: ", string(v))
			}
		}
		break
	//case "aes":
	//	{
	//	}
	//	break
	//case "des":
	//	{
	//	}
	//	break
	//case "hmac":
	//	{
	//	}
	//	break
	//case "rsa":
	//	{
	//	}
	//	break
	default:
		fmt.Println("暂不支持的加密方法, ", action)
	}
}

func init() {
	rootCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")
	signCmd.Flags().StringVarP(&action, "type", "t", "", "选择加密方法")
	signCmd.Flags().BoolVarP(&direct, "direct", "d", false, "加密或者解密,支持部分可逆的算法, 默认加密")

	signCmd.Flags().StringVarP(&sourceFile, "sourceFile", "s", "", "选择计算的加密文件")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
