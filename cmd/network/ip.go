package network

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

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
			ip, _ := getLocalIP()
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

// 获取本机网卡IP
func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	err = &commonError{"获取IP异常"}
	return
}

func getExternal(c *cobra.Command) {
	url, _ := c.Flags().GetString("url")
	resp, err := http.Get(url)
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
