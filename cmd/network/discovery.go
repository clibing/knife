package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	common "github.com/clibing/go-common/pkg/model"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	snet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

// network/discoveryCmd represents the network/discovery command
var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "discovery client list",
	Long: `discovery machine:

show machine by latest sync time.`,
	Run: func(cmd *cobra.Command, args []string) {
		value := requestDiscovery()
		response := &common.Response[[]*Cache]{}
		json.Unmarshal(value, response)

		if response.Code != 200 {
			fmt.Println("发现设备异常,具体: ", response.Message)
			return
		}

		sort.SliceStable(response.Data, func(i, j int) bool {
			return response.Data[i].Metrics.CreateTime.After(response.Data[j].Metrics.CreateTime)
		})

		keyword, _ := cmd.Flags().GetString("keyword")
		skipIpv6, _ := cmd.Flags().GetBool("show-ipv6")

		data := [][]string{}
		for id, v := range response.Data {
			// 网卡名字与当前对应的ip
			interfaces := filterInterfaceIp(v.Metrics.Net.InterfaceStatList, !skipIpv6)
			for _, currentInterface := range interfaces {
				for _, address := range currentInterface.Addrs {
					ipType := "ipv4"
					if !address.IsV4 {
						ipType = "ipv6"
					}
					if len(keyword) > 0 {
						ignorekeywrod := strings.ToLower(keyword)
						if strings.Contains(strings.ToLower(v.Account), ignorekeywrod) || strings.Contains(strings.ToLower(currentInterface.Name), ignorekeywrod) || strings.Contains(strings.ToLower(address.Value), ignorekeywrod) {
							data = append(data, []string{strconv.Itoa(id + 1), v.Account, currentInterface.Name, ipType, address.Value})
						}
					} else {
						data = append(data, []string{strconv.Itoa(id + 1), v.Metrics.CreateTime.Format("2006-01-02 15:04:05.000"), v.Account, currentInterface.Name, ipType, address.Value})
					}
				}
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "LatestTime", "Account", "Device", "Type", "Address"})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render() // Send outpu
	},
}

func filterInterfaceIp(data []snet.InterfaceStatList, skipIpv6 bool) (value []*Addrs) {
	for _, isl := range data {
		for _, is := range isl {
			up, _, _, loopback := checkInterface(is.Flags)
			if !up {
				continue
			}
			if loopback {
				continue
			}
			if strings.HasPrefix(is.Name, "docker") || strings.HasPrefix(is.Name, "br-") || strings.HasPrefix(is.Name, "veth") || strings.HasPrefix(is.Name, "utun") {
				continue
			}
			if len(is.Addrs) == 0 {
				continue
			}

			ips := checkAddrs(is.Addrs, skipIpv6)
			if len(ips) == 0 {
				continue
			}

			value = append(value, &Addrs{
				Name:  is.Name,
				Mac:   is.HardwareAddr,
				Addrs: ips,
			})
		}
	}
	return
}

func checkInterface(flags []string) (up, broadcase, multicase, loopback bool) {
	if len(flags) == 0 {
		return
	}

	for _, v := range flags {
		if strings.Compare("up", v) == 0 {
			up = true
		} else if strings.Compare("broadcase", v) == 0 {
			broadcase = true
		} else if strings.Compare("multicase", v) == 0 {
			multicase = true
		} else if strings.Compare("loopback", v) == 0 {
			loopback = true
		}
	}
	return
}

func checkAddrs(addrs snet.InterfaceAddrList, skipIpv6 bool) (value []*Addr) {
	if len(addrs) == 0 {
		return
	}

	for _, v := range addrs {
		ip, _, e := net.ParseCIDR(v.Addr)
		if e != nil {
			continue
		}

		if p4 := ip.To4(); len(p4) == net.IPv4len {
			value = append(value, &Addr{
				IsV4:  true,
				Value: ip.String(),
			})
			continue
		} else if skipIpv6 {
			continue
		} else {
			value = append(value, &Addr{
				IsV4:  false,
				Value: ip.String(),
			})
		}

	}
	return
}

func init() {
	discoveryCmd.Flags().StringP("keyword", "k", "", "过滤内容，支持主机名字、ip地址、mac地址")
	discoveryCmd.Flags().BoolP("hidden-ipv4", "4", false, "隐藏ipv4地址")
	discoveryCmd.Flags().BoolP("show-ipv6", "6", false, "显示ipv6地址")
}

func NewDiscoveryCmd() *cobra.Command {
	return discoveryCmd
}

func requestDiscovery() (value []byte) {
	url := "https://discovery.linuxcrypt.cn/api/show"

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

	value, _ = io.ReadAll(resp.Body)
	return

}

type Addrs struct {
	Name  string
	Mac   string
	Addrs []*Addr
}

// addr地址
type Addr struct {
	IsV4  bool
	Value string
}

type Cache struct {
	Account string   `json:"account"`
	Metrics *Metrics `json:"metrics"`
}

type Metrics struct {
	CreateTime time.Time `json:"create_time,omitempty"` // 上报事件
	Cpu        *Cpu      `json:"cpu,omitempty"`         //  cpu
	Disk       *Disk     `json:"disk,omitempty"`        // 磁盘
	Host       *Host     `json:"host,omitempty"`        // 主机
	Load       []*Load   `json:"load,omitempty"`        // 负载
	Memory     *Memory   `json:"memory,omitempty"`      // 内存
	Net        *Net      `json:"net,omitempty"`         // 网络配置
	Process    *Process  `json:"process,omitempty"`     // 当前进程
	Hysteria   *Hysteria `json:"hysteria,omitempty"`    // 代理信息
}

type Cpu struct {
	InfoStat []cpu.InfoStat  `json:"info_stat,omitempty"`
	TimeStat []cpu.TimesStat `json:"time_stat,omitempty"`
}

type Disk struct {
	IOCountersStat map[string]disk.IOCountersStat `json:"iocounters_stat,omitempty"`
	PartitionStat  []disk.PartitionStat           `json:"partition_stat,omitempty"`
	UsageStat      *disk.UsageStat                `json:"usage_stat,omitempty"`
}

type Host struct {
	InfoStat *host.InfoStat `json:"info_stat,omitempty"`
}

type Load struct {
	AvgStat *load.AvgStat `json:"avg_stat,omitempty"` // 当前负载
}

type Memory struct {
	SwapDevice        [][]*mem.SwapDevice      `json:"swap_device,omitempty"`         // 交换分区 可能不存在
	SwapMemoryStat    []*mem.SwapMemoryStat    `json:"swap_memory_stat,omitempty"`    // 交换分区 内存统计
	VirtualMemoryStat []*mem.VirtualMemoryStat `json:"virtual_memory_stat,omitempty"` // 系统内容 统计
}

type Net struct {
	InterfaceStatList []snet.InterfaceStatList           `json:"interface_stat_list,omitempty"` // 网卡 mac地址、获取的ip、mtu、
	IOCountersStat    [][]snet.IOCountersStat            `json:"iocounters,omitempty"`          // 网络io
	ConnectionStat    []map[string][]snet.ConnectionStat `json:"connection,omitempty"`          // 查看网卡连接统计，根据Kind类型过滤
	ConntrackStat     [][]snet.ConntrackStat             `json:"conntrack,omitempty"`           //
}

type Process struct {
	Process        process.Process `json:"process,omitempty"`
	Pid            string          `json:"pid,omitempty"`
	Name           string          `json:"name,omitempty"`
	Status         string          `json:"status,omitempty"` // R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock
	Running        bool            `json:"running,omitempty"`
	CreateTime     int64           `json:"createTime,omitempty"`
	MemoryPercent  float32
	CPUPercent     float64
	Groups         []int32
	Cmdline        string
	Exe            string   `json:"exe,omitempty"`     // 可执行文件的完整路径
	Cwd            string   `json:"cwd,omitempty"`     // 工作目录
	Environ        []string `json:"environ,omitempty"` // 当前process使用的env
	MemoryInfoStat *process.MemoryInfoStat
	TimesStat      cpu.TimesStat
	NumThreads     int32                   `json:"num_threads,omitempty"`     // 当前使用多少个threads
	IOCountersStat *process.IOCountersStat `json:"iocounters_stat,omitempty"` // 进程的io统计
}

/**
 * proxy
 */
type Hysteria struct {
	Interface string   `json:"interface,omitempty"` // 网卡名字
	Proxy     []string `json:"proxy,omitempty"`     // 代理url ip:port
}
