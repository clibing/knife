/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"github.com/go-basic/uuid"
)

var (
	needChar, needNumber, needPunctuation bool
	randomLen, randomTimes int
	lettersNumber = "0123456789"
	lettersChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lettersPunctuation = "!@#$%^&*"
)


// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "随机生成指定的内容",
	Long: `生成执行随机数:

1. 生成UUID 默认
2. 生成指定长度的纯数字
3. 生成一个强度比较高的密码
等等
.`,
	Run: func(_ *cobra.Command, _ []string) {
		if !needChar && !needNumber && !needPunctuation {
			for t := 1; t <= randomTimes; t++ {
			    uuid := uuid.New()
                fmt.Println(uuid)
			}
		} else {
			valueLen := randomLen
			result := make([]rune, valueLen)
			r:=rand.New(rand.NewSource(time.Now().UnixNano()))
			randomBase := ""
			if needChar {
				randomBase = randomBase + lettersChar
			}
			if needNumber {
				randomBase = randomBase + lettersNumber
			}
			if needPunctuation {
				randomBase = randomBase + lettersPunctuation
			}
			letters := []rune(randomBase)
			tmp := len(randomBase) - 1
			if randomTimes <= 0 {
				randomTimes = 1
			}

			for t := 1; t <= randomTimes; t++ {
				for i := range result {
					result[i] = letters[r.Intn(tmp)]
				}
				fmt.Println(string(result))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// randomCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	randomCmd.Flags().BoolVarP(&needChar, "needChar", "c", false, "是否需要字符串, 默认不需要")
	randomCmd.Flags().BoolVarP(&needNumber, "needNumber", "n", false, "是否需要数字, 默认不需要")
	randomCmd.Flags().BoolVarP(&needPunctuation, "needPunctuation", "p", false, "是否需要标点符号, 默认不需要")

	randomCmd.Flags().IntVarP(&randomLen, "randomLen", "l", 6, "生成随机数的长度, 默认长度6个字符")
	randomCmd.Flags().IntVarP(&randomTimes, "randomTimes", "t", 1, "生成随机数的个数，默认1个")
}
