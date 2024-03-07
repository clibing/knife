package common

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-basic/uuid"
	"github.com/spf13/cobra"
)

const formatString = "b2eff38a-8832-99de-9330-5fafed0ffacd"
const inputString = "b2eff38a883299de93305fafed0ffacd"

var (
	needChar, needNumber, needPunctuation bool
	randomLen, randomTimes                int
	lettersNumber                         = "0123456789"
	lettersChar                           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lettersPunctuation                    = "!@#$%^&*"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "随机生成指定的内容",
	Long: `生成执行随机数:

1. 生成UUID 默认
2. 生成指定长度的纯数字
3. 生成一个强度比较高的密码
等等
.`,
	Run: func(c *cobra.Command, input []string) {

		ftu, _ := c.Flags().GetBool("format-to-uuid")
		if ftu {
			formatSize := len(formatString)
			inputSize := len(inputString)

			for _, value := range input {
				currentSize := len(value)
				if inputSize == currentSize {
					fmt.Printf("%s-%s-%s-%s-%s\n", value[0:7], value[8:12], value[12:16], value[16:20], value[20:])
				} else if currentSize == formatSize {
					fmt.Println(value)
				}
			}
			return
		}

		if !needChar && !needNumber && !needPunctuation {
			for t := 1; t <= randomTimes; t++ {
				uuid := uuid.New()
				fmt.Println(uuid)
			}
		} else {
			valueLen := randomLen
			result := make([]rune, valueLen)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomBase := ""
			if needChar {
				randomBase = randomBase + lettersChar
			}
			if needNumber {
				randomBase = randomBase + lettersNumber
			}
			if needPunctuation {
				randomBase = randomBase + lettersPunctuation
			}
			letters := []rune(randomBase)
			tmp := len(randomBase) - 1
			if randomTimes <= 0 {
				randomTimes = 1
			}

			for t := 1; t <= randomTimes; t++ {
				for i := range result {
					result[i] = letters[r.Intn(tmp)]
				}
				fmt.Println(string(result))
			}
		}
	},
}

func init() {
	randomCmd.Flags().BoolVarP(&needChar, "needChar", "c", false, "是否需要字符串, 默认不需要")
	randomCmd.Flags().BoolVarP(&needNumber, "needNumber", "n", false, "是否需要数字, 默认不需要")
	randomCmd.Flags().BoolVarP(&needPunctuation, "needPunctuation", "p", false, "是否需要标点符号, 默认不需要")

	randomCmd.Flags().IntVarP(&randomLen, "randomLen", "l", 6, "生成随机数的长度, 默认长度6个字符")
	randomCmd.Flags().IntVarP(&randomTimes, "randomTimes", "t", 1, "生成随机数的个数，默认1个")

	randomCmd.Flags().BoolP("format-to-uuid", "f", false, "将输入的内容格式化为uuid格式")
}

func NewRandomCmd() *cobra.Command {
	return randomCmd
}
