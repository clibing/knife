/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/clibing/knife/cmd/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "多客户端: http, redis, websocket, download, mqtt",
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

	// 下载器
	clientCmd.AddCommand(client.NewDownloadCmd())

	// client --> mqtt
	clientCmd.AddCommand(client.NewMqttCmd())

	rootCmd.AddCommand(clientCmd)
}
