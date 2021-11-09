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
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var path string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装",
	Long: `将当前可执行的二进制程序，安装到系统指定目录下，默认/usr/local/bin:

1. 默认安装
knife install 
2. 指定目录安装
knife install -p /usr/local/bin
.`,
	Run: func(cmd *cobra.Command, args []string) {
		binPath, err := exec.LookPath(os.Args[0])
		if err != nil {
			fmt.Errorf("failed to get bin file info: %s: %s", os.Args[0], err)
			return
		}

		currentFile, err := os.Open(binPath)
		if err != nil {
			fmt.Errorf("failed to get bin file info: %s: %s", binPath, err)
			return
		}
		defer func() { _ = currentFile.Close() }()

		installFile, err := os.OpenFile(filepath.Join(path, "knife"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Errorf("failed to create bin file: %s: %s", filepath.Join(path, "knife"), err)
			return
		}
		defer func() { _ = installFile.Close() }()

		_, err = io.Copy(installFile, currentFile)
		if err != nil {
			fmt.Errorf("failed to copy file: %s: %s", filepath.Join(path, "knife"), err)
			return
		}
		fmt.Println("install success")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")
	installCmd.PersistentFlags().StringVarP(&path, "path", "p", "/usr/local/bin", "安装目录，window需要指定目录")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
