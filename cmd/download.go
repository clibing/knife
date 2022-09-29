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
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var storagePath string
var sourceUrl string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "下载器",
	Long: `下载器:

1. 快速下载, 直接接url
knife download https://clibing.com/download.demo.tar.gz

2. 指定下载目录并指定下载的资源地址
knife download -p ./ -s https://clibing.com/download.demo.tar.gz
.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(storagePath) == 0 && len(sourceUrl) == 0 {
			for _, url := range args {
				Download(url, "./")
			}
		} else {
			Download(sourceUrl, storagePath)
		}
	},
}

func Download(source, path string) {
	req, _ := http.NewRequest("GET", source, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("请求下载失败", err)
		return
	}
	defer resp.Body.Close()

	var filename string
	// 从协议头 提取文件的名字
	header := resp.Header["Content-Disposition"]
	if len(header) > 0 {
		//`attachment;filename="foo.png"`
		_, params, err := mime.ParseMediaType(header[0])
		if err == nil {
			fmt.Println("解析服务端http协议头异常")
		} else {
			filename = params["filename"]
		}
	}
	if len(filename) == 0 {
		filename = "tmp"
		u, err := url.Parse(source)
		if err == nil {
			t := strings.Split(u.Path, "/")
			filename = t[len(t)-1]
		}
	}

	f, _ := os.OpenFile(fmt.Sprintf("%s/%s", path, filename), os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"进度",
	)

	io.Copy(io.MultiWriter(f, bar), resp.Body)
}
func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	downloadCmd.Flags().StringVarP(&storagePath, "storagePath", "p", "", "存储目录")
	downloadCmd.Flags().StringVarP(&sourceUrl, "sourceUrl", "s", "", "资源地址")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
