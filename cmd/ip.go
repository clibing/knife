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
	"io/ioutil"
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
	Run: func(cmd *cobra.Command, args []string) {
		if external {
			getExternal()
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

func getExternal() {
	resp, err := http.Get("https://ipw.cn/api/ip/myip")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("external ip: ", string(s))
	// io.Copy(os.Stdout, resp.Body)
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(ipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipCmd.PersistentFlags().String("foo", "", "A help for foo")
	ipCmd.Flags().BoolVarP(&external, "external", "e", false, "获取出口ip，默认获取本机ip")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
