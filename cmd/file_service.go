/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	port       int
	staticPath string
	// values     []string
)

// staticCmd represents the static command
var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "文件服务。",
	Long: `启用本地静态资源服务:

新装系统后，安装所需软件的时候，每次都需要移动硬盘、U盘或者scp等拷贝资源到目标机器。
一般情况都有一台闲置的电脑, 被安装的电脑在安装机器的期间, 可以使用闲置的机器可以去官网现在所需最新的软件安装包。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取端口
		value := fmt.Sprintf(":%d", port)
		// 静态资源的目录
		fs := http.FileServer(http.Dir(staticPath))
		// http 处理器
		http.Handle("/", http.StripPrefix("/", fs))
		// 建立监听
		listener, err := net.Listen("tcp", value)
		if err != nil {
			fmt.Println("建立监听异常, ", err.Error())
			return
		}

		fmt.Printf("服务启动中: http://0.0.0.0:%d \n", listener.Addr().(*net.TCPAddr).Port)
		err = http.Serve(listener, nil)
		if err != nil {
			fmt.Println("服务启动失败，请检查, ", err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(staticCmd)

	// Here you will define your flags and configuration settings.

	staticCmd.Flags().StringVarP(&staticPath, "path", "p", "", "静态资源目录, 默认为当前目录")
	staticCmd.Flags().IntVarP(&port, "port", "", 0, "端口, 默认会随机")
	// staticCmd.Flags().StringSliceVarP(&values, "values", "v", nil, "parameters")

	// staticCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
