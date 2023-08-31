package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
