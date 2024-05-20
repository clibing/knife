package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// httpCmd represents the static command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "http服务。",
	Long: `启用http服务:

knife server http --port=8080 --uri=/ 返回 json: '{\"code\":200, \"ts\":\"now()\"}' 
`,
	Run: func(cmd *cobra.Command, args []string) {

		port, _ := cmd.Flags().GetInt("port")
		uri, _ := cmd.Flags().GetString("uri")

		http.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("{\"code\":200, \"message\":\"%d\"}", time.Now().UnixMilli())))
		})

		// 获取端口
		value := fmt.Sprintf(":%d", port)
		// 建立监听
		listener, err := net.Listen("tcp", value)
		if err != nil {
			fmt.Println("建立监听异常, ", err.Error())
			return
		}

		// 获取 本地ip地址
		port = listener.Addr().(*net.TCPAddr).Port
		fmt.Printf("服务启动中: http://0.0.0.0:%d uri: %s\n", port, uri)
		err = http.Serve(listener, nil)
		if err != nil {
			fmt.Println("服务启动失败，请检查, ", err.Error())
			return
		}

	},
}

func init() {
	httpCmd.Flags().IntP("port", "", 0, "端口, 默认会随机")
	httpCmd.Flags().StringP("uri", "", "/", "uri处理")
}

func NewHttpServer() *cobra.Command {
	return httpCmd
}
