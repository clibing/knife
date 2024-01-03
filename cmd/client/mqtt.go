package client

import "github.com/spf13/cobra"

// mqttCmd represents the random command
var mqttCmd = &cobra.Command{Use: "mqtt",
	Short: "mqtt客户端",
	Long: `mqtt客户端:

1. 订阅
2. 发布
.`,
	Run: func(cmd *cobra.Command, arg []string) {

	},
}

var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "发布",
	Long: `发布消息:

1. 发布消息
  knife client mqtt pub -t /Offline -m 'message'		

		`,
	Run: func(cmd *cobra.Command, arg []string) {

	},
}

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "订阅",
	Long: `订阅消息:

1. 发布消息
  knife client mqtt sub -t /Offline

		`,
	Run: func(cmd *cobra.Command, arg []string) {

	},
}

func init() {
	// mqttCmd.Flags().IntVarP(&randomTimes, "randomTimes", "t", 1, "生成随机数的个数，默认1个")
}

func NewMqttCmd() *cobra.Command {
	mqttCmd.AddCommand(pubCmd)
	mqttCmd.AddCommand(subCmd)
	return mqttCmd
}
