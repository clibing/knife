package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/clibing/knife/cmd/system"
	"github.com/spf13/cobra"
)

var path string

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "卸载",
	Long: `将当前可执行的二进制程序，从到系统指定目录下移除，默认目录/usr/local/bin:

1. 默认卸载
knife uninstall 
2. 指定目录卸载
knife uninstall -p /usr/local/bin.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("uninstall called")
	},
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装",
	Long: `将当前可执行的二进制程序，安装到系统指定目录下，默认/usr/local/bin:

1. 默认安装
knife install 
2. 指定目录安装
knife install -p /usr/local/bin
.`,
	Run: func(_ *cobra.Command, _ []string) {
		binPath, err := exec.LookPath(os.Args[0])
		if err != nil {
			fmt.Printf("failed to get bin file info: %s: %s", os.Args[0], err)
			return
		}

		currentFile, err := os.Open(binPath)
		if err != nil {
			fmt.Printf("failed to get bin file info: %s: %s", binPath, err)
			return
		}
		defer func() { _ = currentFile.Close() }()

		installFile, err := os.OpenFile(filepath.Join(path, "knife"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Printf("failed to create bin file: %s: %s", filepath.Join(path, "knife"), err)
			return
		}
		defer func() { _ = installFile.Close() }()

		_, err = io.Copy(installFile, currentFile)
		if err != nil {
			fmt.Printf("failed to copy file: %s: %s", filepath.Join(path, "knife"), err)
			return
		}
		fmt.Println("install success")
	},
}

var systemCmd = &cobra.Command{
	Use:     "system",
	Aliases: []string{"sys"},
	Short:   `系统工具: arch, monitor, upgrade, maven`,
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看当前版本号",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("version: %s\n  build: %s\n commit: %s\n author: clibing (wmsjhappy@gmail.com)\nCopyright (c) %s clibing, All rights reserved.",
			version, buildDate, commitId, time.Now().Format("2006"))
	},
}

func init() {
	systemCmd.AddCommand(
		system.NewArchCmd(),
		system.NewMonitorCmd(),
		system.NewUpgradeCmd(),
		system.NewMavenCmd(),
		system.NewCronCmd(),
	)

	// 转换器
	rootCmd.AddCommand(systemCmd)

	// 安装 cmd
	installCmd.PersistentFlags().StringVarP(&path, "path", "p", "/usr/local/bin", "安装目录，window需要指定目录")
	rootCmd.AddCommand(installCmd)

	// 下载cmd
	rootCmd.AddCommand(uninstallCmd)

	// 系统版本
	rootCmd.AddCommand(versionCmd)
}
