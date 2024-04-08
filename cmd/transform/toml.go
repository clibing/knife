package transform

import "github.com/spf13/cobra"


// tomlCmd format & edit
var tomlCmd = &cobra.Command{
	Use:   "toml",
	Short: "toml文件格式化、编辑等",
	Example: `

`,
	Run: func(c *cobra.Command, args []string) {
	},
}

func init() {
	tomlCmd.Flags().StringP("input", "i", "\t", "json格式化缩进标记，默认制表符")
	tomlCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "json格式化前缀")
	tomlCmd.Flags().IntVarP(&convert, "convert", "c", 0, "转换模式\n0: json格式化\n1: xml to json\n2: json to xml(建议使用struct Tag)\n3: json to yml\n4: yml to json, 默认为0美化")
}

func NewTomlCommand() *cobra.Command {
	return convertCmd
}
