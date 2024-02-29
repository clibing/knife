package system

import (
	"fmt"

	"github.com/clibing/knife/internal/utils"
	"github.com/spf13/cobra"
)

// beautifyCmd represents the cpu command
var beautifyCmd = &cobra.Command{
	Use:   "beautify",
	Short: "字节美化工具, 反向计算字节",
	Long: `example: 

将字节转换为美化的+单位
knife system beautify -b 1026
	
	`,
	Run: func(c *cobra.Command, _ []string) {
		size, _ := c.Flags().GetInt64("byte")
		name := utils.BeautifyValue(size)
		fmt.Println(name)
	},
}

func init() {
	beautifyCmd.Flags().Int64P("byte", "b", 0, "字节")
}

func NewBeautifyCmd() *cobra.Command {
	return beautifyCmd
}
