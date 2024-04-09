package network

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/clibing/knife/internal/utils"
	"github.com/spf13/cobra"
)

var external bool

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "获取ip地址",
	Long: `可以获取本机的ip地址和本机的出口ip:

1. 获取本机ip
knife ip 
2. 获取出口ip
knife ip -e

本机IP地址使用的时候比较多，出口ip一般可以做一个端口处理转发时使用的比较多，本功能较为简单.`,
	Run: func(c *cobra.Command, _ []string) {
		if external {
			getExternal(c)
		} else {
			ip, _ := utils.GetLocalIp()
			fmt.Println("local ip: ", ip)
		}
	},
}

type commonError struct {
	title string
}

func (ce *commonError) Error() string {
	return ce.title
}

func getExternal(c *cobra.Command) {
	url, _ := c.Flags().GetString("url")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "clibing/knife")
	resp, err := client.Do(req)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()

	s, _ := io.ReadAll(resp.Body)
	fmt.Println("external ip: ", string(s))
	// io.Copy(os.Stdout, resp.Body)
	os.Exit(0)
}

func init() {
	ipCmd.Flags().StringP("url", "u", "https://discovery.linuxcrypt.cn/api/ip", "请求url")
	ipCmd.Flags().BoolVarP(&external, "external", "e", false, "获取出口ip，默认获取本机ip")
}

func NewIpCmd() *cobra.Command {
	return ipCmd
}
