package sign

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var md5Cmd = &cobra.Command{
	Use:   "md5",
	Short: `md5签名计算`,
	Run: func(c *cobra.Command, args []string) {
		input, _ := c.Flags().GetString("input")
		if len(args) == 0 && len(input) == 0 {
			c.Help()
		}
		if len(input) != 0 {
			value, err := os.ReadFile(sourceFile)
			if err == nil {
				fmt.Println("sign file error, ", err)
			} else {
				md5sum(value)
			}
		}
		if len(args) > 0 {
			for _, content := range args {
				md5sum([]byte(content))
			}
		}
	},
}

func md5sum(content []byte) {
	h := md5.New()
	h.Write([]byte(content))
	value := hex.EncodeToString(h.Sum(nil))
	fmt.Println("source :", string(content))
	fmt.Println("md5(32):", value)
	fmt.Println("       :", strings.ToUpper(value))
	fmt.Println("md5(16):", value[8:24])
	fmt.Println("       :", strings.ToUpper(value[8:24]))
}
func init() {
	md5Cmd.Flags().StringP("input", "i", "", "输入待验证的文件")
	// md5Cmd.Flags().StringP("input", "i", "", "输入待验证的文件")
}

func NewMd5Cmd() *cobra.Command {
	return md5Cmd
}
