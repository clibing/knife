package sign

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var sha1Cmd = &cobra.Command{
	Use:   "sha1",
	Short: `sha1签名计算`,
	Run: func(c *cobra.Command, args []string) {
		input, _ := c.Flags().GetString("input")
		if len(args) == 0 && len(input) == 0 {
			c.Help()
		}
		if len(input) != 0 {
			value, err := os.ReadFile(input)
			if err != nil {
				fmt.Println("sign file error, ", err)
			} else {
				sha1sum(value)
			}
		}
		if len(args) > 0 {
			for _, content := range args {
				sha1sum([]byte(content))
			}
		}
	},
}

func sha1sum(content []byte) {
	s := sha1.New()
	s.Write(content)
	value := hex.EncodeToString(s.Sum(nil))
	fmt.Println("source :", strings.Replace(string(content), "\n", "", -1))
	fmt.Println("output :", value)
	fmt.Println("       :", strings.ToUpper(value))

}
func init() {
	sha1Cmd.Flags().StringP("input", "i", "", "输入待验证的文件")
}

func NewSha1Cmd() *cobra.Command {
	return sha1Cmd
}
