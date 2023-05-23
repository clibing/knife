/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
    str, ver, o, arch string
)
// archCmd represents the arch command
var archCmd = &cobra.Command{
	Use:   "arch",
	Short: "get current system kernal, machine.",
	Long: `get current os env:

默认替换os、arch关键字, 注意输入的-s的值需要使用单引号, 不能使用双信号。

knife arch -s 'https://go.dev/dl/go${version}.${os}-${arch}.tar.gz' 

knife arch -s 'https://go.dev/dl/go${version}.${os}-${arch}.tar.gz' -o linux -a 386 -v 1.20.4

使用环境变量
set VV=99
knife arch -s 'https://go.dev/dl'"$VV"'/go${version}.${os}-${arch}.tar.gz' -o linux -a 386 -v 1.20.4

$GOOS	$GOARCH
darwin	386
darwin	amd64
darwin	arm
darwin	arm64
dragonfly	amd64
freebsd	386
freebsd	amd64
freebsd	arm
linux	386
linux	amd64
linux	arm
linux	arm64
linux	ppc64
linux	ppc64le
netbsd	386
netbsd	amd64
netbsd	arm
openbsd	386
openbsd	amd64
openbsd	arm
plan9	386
plan9	amd64
solaris	amd64
windows	386
windows	amd64

current os arch, kernal.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(str) == 0 {
			fmt.Printf("%s %s\n", runtime.GOOS, runtime.GOARCH)
		}

		result := str
		if len(ver) > 0 {
			result = strings.ReplaceAll(result, "${version}", ver)
		}

		if len(o) > 0 {
			result = strings.ReplaceAll(result, "${os}", o)
		} else {
			result = strings.ReplaceAll(result, "${os}", runtime.GOOS)
		}

		if len(arch) > 0 {
			result = strings.ReplaceAll(result, "${arch}", arch)
		}else {
			result = strings.ReplaceAll(result, "${arch}", runtime.GOARCH)
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(archCmd)

	// Here you will define your flags and configuration settings.
    archCmd.PersistentFlags().StringVarP(&str, "str", "s", "", "需要替换的URL")
    archCmd.PersistentFlags().StringVarP(&ver, "version", "v", "", "需要替换的URL")
    archCmd.PersistentFlags().StringVarP(&o, "os", "o", "", "当前系统运行的操作系统")
    archCmd.PersistentFlags().StringVarP(&arch, "arch", "a", "", "当前系统运行的架构")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// archCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// archCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
