package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

/**
 * command: 当前执行的指令
 * err: 当前错误
 * debug: debug 错误信息
 * product: 生产环境 提示的信息
 */
func message(c *cobra.Command, develop_msg, product_msg string, entry bool, value ...interface{}) {
	debug, e := c.Flags().GetBool(fmt.Sprintf("%s-debug", c.Use))
	if e != nil {
		debug = false
	}
	if debug {
		fmt.Printf(setEntry(develop_msg, entry), value...)
		return
	}
	if len(product_msg) > 0 {
		fmt.Printf(setEntry(product_msg, entry), value...)
	}
}

func setEntry(source string, entry bool) string {
	if entry {
		return source + "\n"
	}
	return source

}
