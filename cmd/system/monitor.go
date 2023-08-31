/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
package system

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

var (
	times int
)

// monitorCmd represents the cpu command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "检查当前系统cpu使用率，内存使用率",
	Run: func(_ *cobra.Command, _ []string) {
		var i int
		if times > 100 {
			fmt.Println("times超出限制，默认为100")
			times = 99
		}
		for i = 0; i < times; i++ {
			v, _ := mem.VirtualMemory()
			fmt.Printf("- %02d 当前cpu的温度为: %.2f℃, 内存: %.2fG(%v字节), 剩余内存: %.2fG(%v字节), 使用率: %.2f%%\n", i+1, termperature(), float64(v.Total)/1024/1024/1024, v.Total, float64(v.Free)/1024/1024/1024, v.Free, v.UsedPercent)
		}
	},
}

func termperature() (val float64) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("获取失败")
	}
	val = percent[0]
	return
}
func init() {
	monitorCmd.Flags().IntVarP(&times, "times", "t", 3, "检查cpu当前温度，需要检查多少次，间隔为1秒")
}

func NewMonitorCmd() *cobra.Command {
	return monitorCmd
}
