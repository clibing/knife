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
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var (
	// 下载地址
	download = ""
	// 升级的版本
	upgrade = ""
	// 二进制的位置
	bin = ""
)

const downloadURL = "https://github.com/clibing/knife/releases/download/%s/knife_%s_%s"

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "客户端升级",
	Long: `通过指定版本进行快速升级:

knife upgrade -u "0.0.6" -b /usr/local/bin/
.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(download) == 0 && len(upgrade) == 0 {
			fmt.Println("参数不能为空")
			return
		}

		target := ""
		if len(download) != 0 {
			target = download
		} else {
			arch := runtime.GOARCH
			os := runtime.GOOS
			target = fmt.Sprintf(downloadURL, upgrade, os, arch)
		}

		fmt.Println("下载的地址: ", target)

		file := filepath.Join(path, bin)

		// backup
		if fileExist(file) {
			if e :=autoBackup(file); e != nil {
				fmt.Println(e.Error())
				return
			}
		}

		res, err := http.Get(target)
		defer func() { _ = res.Close }()

		if err != nil {
			fmt.Println("下载失败: ", target)
			return
		}

		pf, err := os.Create(file)
		defer func() { _ = pf.Close }()
		if err != nil {
			fmt.Println("下载失败: ", target)
			return
		}

		io.Copy(pf, res.Body)
		os.Chmod(file, 0777)
	},
}

func autoBackup(file string) (error){
	bak, err := os.Open(file)
	if err != nil {
		fmt.Println("打开现有文件异常", err)
		return err
	}
	defer func() { _ = bak.Close() }()

	suffix := version
	if len(version) == 0 {
		suffix = strconv.Itoa(int(time.Now().UnixMilli() / 1000))
	}
	bakF := fmt.Sprintf("%s.%s", file, suffix)
	pp, err := os.Create(bakF)
	if err != nil {
		fmt.Println("创建备份异常: ", err)
		return err
	}
	defer func() { _ = pp.Close() }()
	io.Copy(pp, bak)
	fmt.Println("备份数据成功, ", bakF)
	return nil
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	upgradeCmd.PersistentFlags().StringVarP(&download, "download", "d", "", "下载url")
	upgradeCmd.PersistentFlags().StringVarP(&upgrade, "upgrade", "u", "", "升级的tag，例如v0.0.1")
	upgradeCmd.PersistentFlags().StringVarP(&path, "path", "p", "/usr/local/bin", "二进制的保存目录, 默认为 /usr/local/bin")
	upgradeCmd.PersistentFlags().StringVarP(&bin, "bin", "b", "knife", "二进制安装的名字, 默认为 knife")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
