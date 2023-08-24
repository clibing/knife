/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "client",
		Long:  `多应用客户端`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("client called")
		},
	}
	clientType, clientTarget string
)

func init() {

	clientCmd.AddCommand(&cobra.Command{
		Use:   "http",
		Short: "http",
		Long:  `http client`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("http client called")
		},
	})

	rootCmd.AddCommand(clientCmd)

	// clientCmd.Flags().StringVarP(&clientType, "client-type", "c", "", "选择客户端类型： http, socks5, websocket")
	clientCmd.Flags().StringVarP(&clientType, "client-type", "c", "", `选择客户端类型：
    http: 请求url, 类curl。
    socks5: proxy,通过http请求测试。
    websocket: ws客户端。`)
	clientCmd.Flags().StringVarP(&clientTarget, "client-target", "t", "", "请求目标， ip地址或者url")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
