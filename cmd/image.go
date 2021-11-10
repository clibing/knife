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
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var (
	imageFile, base64File  string
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "图片转base64",
	Long: `将图片转换成base64， 也可以将base64生成图片:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if direct {
			// base64 -> image
			ddd, _ := base64.StdEncoding.DecodeString(base64File)  //成图片文件并把文件写入到buffer
			err := ioutil.WriteFile(imageFile, ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
			if err != nil {
				fmt.Errorf("base64 to image error, %s", err)
			}
		} else {
			// image -> base64
			file, err := ioutil.ReadFile(imageFile)
			if err != nil {
				fmt.Errorf("read file error, %s", err)
			}
			result := base64.StdEncoding.EncodeToString(file)
			ioutil.WriteFile(base64File, []byte(result), 0664)
		}
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageCmd.PersistentFlags().String("foo", "", "A help for foo")
	imageCmd.Flags().StringVarP(&imageFile, "imageFile", "i", "", "图片文件")
	imageCmd.Flags().StringVarP(&base64File, "base64File", "b", "", "base64存储的文件")
	imageCmd.Flags().BoolVarP(&direct, "direct", "d", false, "base64生成图片，默认false， 图片生成base64")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
