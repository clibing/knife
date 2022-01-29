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
package cmd

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/spf13/cobra"
)

var (
	times int
)

// cpuCmd represents the cpu command
var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "cpu temperature",
	Run: func(cmd *cobra.Command, args []string) {
		var i int
		if times > 100 {
			times = 99
		}
		for i = 0; i < times; i++ {
			fmt.Printf("%02d 当前cpu的温度为: %.2f℃\n", i+1, termperature())
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
	rootCmd.AddCommand(cpuCmd)

	// Here you will define your flags and configuration settings.
	cpuCmd.Flags().IntVarP(&times, "times", "t", 3, "检查cpu当前温度，需要检查多少次，间隔为1秒")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cpuCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cpuCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
