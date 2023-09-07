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
package backup

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// duration int 已经定义 time.go
	requestDuration int
	requestURL      string
	// url string 已经定义 url.go

)

// abCmd represents the ab command
var abCmd = &cobra.Command{
	Use:   "ab",
	Short: "ab测试接口",
	Long: `ab测试接口:

.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("ab called")
	},
}

// abCmd represents the ab command
var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "bench测试接口",
	Long: `bench测试接口:

.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("sub command bench called")
	},
}

func init() {
	abCmd.AddCommand(benchCmd)
	// rootCmd.AddCommand(abCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	abCmd.PersistentFlags().IntVarP(&requestDuration, "duration", "d", 0, "请求的持续时间")
	abCmd.PersistentFlags().StringVarP(&requestURL, "url", "u", "", "请求的URL")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// abCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
