package network

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

var (
	mac, nic string
)

const UDPPort = 9

// wolCmd represents the wol command
var wolCmd = &cobra.Command{
	Use:   "wol",
	Short: "Wake on Lan",
	Long: `局域网唤醒主机:

原理：
    16进制： 0xFFFFFFFFFF+16次重复的目标mac地址
   
条件:
    1. 被唤醒的主机，主板支持Wake ON Lan并开启
    2. 被唤醒的主机，需要与执行唤醒动作的主机在同一局域网
    3. 被唤醒的主机是正常关机并接入有线。强制关机的不可用。

mac地址格式：
    1. 11:22:33:44:55:66
    2. 11-22-33-44-55-66
    3. 11:22-33:44:55-66
    4. 112233445566

使用：
    knife wol -m 11:22:33:44:55:66
    knife wol -m 11:22:33:44:55:66 -n eth0.`,
	Run: func(c *cobra.Command, _ []string) {
		if mac == "" {
			c.Help()
			return
		}
		fmt.Printf("发送唤醒的网卡地址: %s\n", mac)

		handleMac := strings.ReplaceAll(strings.ReplaceAll(mac, ":", ""), "-", "")
		if len(handleMac) != 12 {
			fmt.Printf("被唤醒的MAC地址不合法: %s\n", mac)
			return
		}

		macHex, err := hex.DecodeString(handleMac)
		if err != nil {
			fmt.Printf("被唤醒的MAC地址不合法: %s\n", mac)
			return
		}

		// 构建唤醒魔包
		var bcast = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
		var buff bytes.Buffer
		buff.Write(bcast)
		for i := 0; i < 16; i++ {
			buff.Write(macHex)
		}

		// 获得唤醒魔包
		data := buff.Bytes()

		sender := net.UDPAddr{}

		currentIp, err := interfaceIPv4(nic)
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			if currentIp != nil {
				sender.IP = currentIp
			}
		}

		target := net.UDPAddr{
			IP:   net.IPv4bcast,
			Port: UDPPort,
		}

		conn, err := net.DialUDP("udp", &sender, &target)
		if err != nil {
			fmt.Println("创建客户端异常", err.Error())
			return
		}
		defer func() {
			_ = conn.Close()
		}()

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("发送唤醒失败", err.Error())
			return
		}
		fmt.Println("指令已经发送， 请稍等.. 主机启动中")
	},
}

func interfaceIPv4(nic string) (net.IP, error) {
	n := strings.TrimSpace(nic)
	if len(n) == 0 {
		return nil, nil
	}

	inter, err := net.InterfaceByName(n)
	if err != nil {
		return nil, fmt.Errorf("获取指定网卡异常, error: %s", err)
	}
	check := inter.Flags & net.FlagUp
	if check == 0 {
		fmt.Println("")
		return nil, errors.New("当前网卡不可用")
	}
	addrs, err := inter.Addrs()
	if err != nil {
		return nil, fmt.Errorf("获取指定网卡的IP异常, error: %s", err)
	}

	var currentIp net.IP
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok {
			if ipv4 := ip.IP.To4(); ipv4 != nil {
				currentIp = ipv4
				break
			}
		}
	}
	if currentIp == nil {
		return nil, fmt.Errorf("获取指定网卡的IP异常, error: %s", err)
	}
	return currentIp, nil
}

func init() {
	wolCmd.PersistentFlags().StringVarP(&mac, "mac", "m", "", "被唤醒的mac地址")
	wolCmd.PersistentFlags().StringVarP(&nic, "nic", "n", "", "发出唤醒的数据包网卡")
}

func NewWolCmd() *cobra.Command {
	return wolCmd
}
