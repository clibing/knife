/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/knife/cmd/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "多客户端: http, redis, websocket",
		Long:  `多应用客户端`,
		Example: `knife client http https://tool.linuxcrypt.cn/checkRemoteIp
knife client ws wss://api.linuxcrypt.cn
		`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	websocketCmd = &cobra.Command{
		Use:     "websocket",
		Short:   "websocket",
		Aliases: []string{"ws"},
		Long:    `websocket client`,
		Run: func(ws *cobra.Command, args []string) {
			fmt.Println("websocket client called")
		},
	}
)

func init() {
	// client -> http request
	httpCmd := client.NewHttpClient()
	clientCmd.AddCommand(httpCmd)

	// client -> redis
	redisCmd := client.NewRedisClient()
	clientCmd.AddCommand(redisCmd)

	websocketCmd.Flags().Bool("websocket-debug", false, "debug, default false")
	// client -> websocket
	clientCmd.AddCommand(websocketCmd)

	rootCmd.AddCommand(clientCmd)

	// clientCmd.Flags().StringVarP(&clientType, "client-type", "c", "", "选择客户端类型： http, socks5, websocket")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
