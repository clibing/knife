package sign

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/clibing/knife/internal/crypto"
	"github.com/spf13/cobra"
)

var aesCmd = &cobra.Command{
	Use:   "aes",
	Short: `aes加密解密`,
	Example: `1. 快速签名计算
knife sign aes -k 1234567890abcedf clibing

2. 解密,自动识别hex和base64
knife sign aes -k 1234567890abcedf -d e0991018dd517965c7ab8311af0885b7
knife sign aes -k 1234567890abcedf -d 4JkQGN1ReWXHq4MRrwiFtw==

3. 使用ECB模式
加密
knife sign aes -k 1234567890abcedf -m ECB clibing
解密
knife sign aes -k 1234567890abcedf -m ECB -d 867a374e5caee0f4c8462a658cc431b7
解密并输出到文件
knife sign aes -k 1234567890abcedf -m ECB -d 867a374e5caee0f4c8462a658cc431b7 -o /tmp/out.txt

4. aes计算文件
knife sign aes -k 1234567890abcedf -i /tmp/input.txt

5. 指定便宜量， 默认偏移量同key
knife sign aes -k 1234567890abcedf --iv 0234567890abcedh clibing`,
	Run: func(c *cobra.Command, args []string) {
		input, _ := c.Flags().GetString("input")
		if len(args) == 0 && len(input) == 0 {
			c.Help()
		}
		if len(input) > 0 {
			value, err := os.ReadFile(input)
			if err != nil {
				fmt.Println("sign file error, ", err)
			} else {
				aseCrypt(c, value)
			}
		}
		if len(args) > 0 {
			for _, content := range args {
				aseCrypt(c, []byte(content))
			}
		}
	},
}

func aseCrypt(c *cobra.Command, content []byte) {
	decrypt, _ := c.Flags().GetBool("decrypt")
	output, _ := c.Flags().GetString("output")
	input, _ := c.Flags().GetString("input")
	show_file := len(input) > 0
	skip_console := len(output) > 0
	mode, _ := c.Flags().GetString("mode")
	key, _ := c.Flags().GetString("key")

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("密钥的长度: 16, 24, 32")
		return
	}

	if len(mode) == 0 {
		mode = "CBC"
	}

	var result []byte
	if !skip_console {
		if show_file {
			fmt.Printf("%14s : %s\n", "source", input)
		} else {
			fmt.Printf("%14s : %s\n", "source", strings.Replace(string(content), "\n", "", -1))
		}
	}

	if decrypt {
		var e error
		if crypto.IsHex(content) {
			content, e = hex.DecodeString(string(content))
		} else {
			content, e = base64.StdEncoding.DecodeString(string(content))
		}
		if e != nil {
			fmt.Println("待解密内容异常", e.Error())
			return
		}
	}

	if !skip_console {
		fmt.Printf("%14s : %s\n", "mode", mode)
	}
	if mode == "CBC" {
		ivInput, _ := c.Flags().GetString("iv")
		iv, e := aesIV(ivInput, key, crypto.BlockSize)
		if e != nil {
			fmt.Printf("%14s : %s\n", "err", e.Error())
			return
		}
		var err error
		if !decrypt {
			result, err = crypto.CBCEncrypt(content, []byte(key), iv)
		} else {
			result, err = crypto.CBCDecrypt(content, []byte(key), iv)
		}
		if err != nil {
			fmt.Printf("%14s : %s\n", "err", err.Error())
			return
		}
	} else if mode == "ECB" {
		if !decrypt {
			result = crypto.ECBEncrypt(content, []byte(key))
		} else {
			result = crypto.ECBDecrypt(content, []byte(key))
		}
	} else if mode == "CFB" {
		var err error
		if !decrypt {
			result, err = crypto.CFBEncrypt(content, []byte(key))
		} else {
			result, err = crypto.CFBDecrypt(content, []byte(key))
		}
		if err != nil {
			fmt.Printf("%14s : %s\n", "err", err.Error())
			return
		}
	} else {
		fmt.Println("Mode not supported: [CBC, ECB, CFB]", mode)
		return
	}

	if !skip_console {
		fmt.Printf("%14s : %s\n", "key", key)
	}

	if decrypt {
		if !skip_console {
			fmt.Printf("%14s : %s\n", "value", string(result))
		} else {
			os.WriteFile(output, result, 0644)
		}
	} else {
		value := hex.EncodeToString(result)
		value2 := base64.StdEncoding.EncodeToString(result)
		if !skip_console {
			fmt.Printf("%14s : %s\n", "value(hex)", value)
			fmt.Printf("%14s : %s\n", "value(base64)", value2)
		} else {
			os.WriteFile(output, []byte(fmt.Sprintf("%14s : %s\n%14s : %s\n", "value(hex)", value, "value(base64)", value2)), 0664)
		}
	}
}

// aesIV 解析 CBC 偏移量：默认同 key（过长则截取），过短报错
func aesIV(input, key string, blockSize int) ([]byte, error) {
	if len(input) == 0 {
		if len(key) > blockSize {
			key = key[:blockSize]
		}
		return []byte(key), nil
	}
	if len(input) < blockSize {
		return nil, fmt.Errorf("偏移量小于块大小, 当前块大小: %d", blockSize)
	}
	if len(input) > blockSize {
		fmt.Printf("偏移量过长，按照块大小(%d)自动截取\n", blockSize)
		input = input[:blockSize]
	}
	return []byte(input), nil
}

func init() {
	aesCmd.Flags().StringP("key", "k", "", "加密的密钥")
	aesCmd.Flags().StringP("input", "i", "", "输入文件")
	aesCmd.Flags().StringP("output", "o", "", "输出的文件")
	aesCmd.Flags().String("iv", "", "当model为CBC时生效，偏移量可用，当为空时使用密钥")
	aesCmd.Flags().StringP("mode", "m", "CBC", "模式: CBC, ECB, CFB")
	aesCmd.Flags().BoolP("decrypt", "d", false, "编码方式，默认encrypt")
}

func NewAesCmd() *cobra.Command {
	return aesCmd
}
