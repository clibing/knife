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
package backup

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var match bool
var expression string

// regCmd represents the reg command
var regCmd = &cobra.Command{
	Use:   "reg",
	Short: "正则表达式",
	Long: `验证正则表达式对内容的匹配或者查找:

knife reg -e "H(.*)d" "HelloWorld message ".`,
	Run: func(_ *cobra.Command, args []string) {
		r, _ := regexp.Compile(expression)
		if match {
			for _, source := range args {
				fmt.Println("match result: ", r.MatchString(source))
			}
		} else {
			for _, source := range args {
				fmt.Println("find result: ", r.FindString(source))
			}
		}
	},
}

func init() {
	// rootCmd.AddCommand(regCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// regCmd.PersistentFlags().String("foo", "", "A help for foo")
	regCmd.Flags().StringVarP(&expression, "expression", "e", "", "正则表达式")
	regCmd.Flags().BoolVarP(&match, "match", "m", false, "是否为match模式, 默认为查找模式")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// regCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
