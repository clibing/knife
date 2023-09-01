package sign

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var sha1Cmd = &cobra.Command{
	Use:   "sha",
	Short: `sha签名计算, bit: 1, 224, 256, 384, 512`,
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
				sha1sum(c, value)
			}
		}
		if len(args) > 0 {
			for _, content := range args {
				sha1sum(c, []byte(content))
			}
		}
	},
}

func sha1sum(c *cobra.Command, content []byte) {
	var handle hash.Hash
	bit, e := c.Flags().GetInt("bit")
	if e != nil {
		fmt.Println("bit not found, (1, 224, 256, 384, 512, 512224, 512256)")
		return
	}
	bits := make([]int, 0)

	if bit == 0 {
		bits = append(bits, 1, 224, 256, 384, 512, 512224, 512256)
	} else if bit == 1 {
		bits = append(bits, 1)
	} else if bit == 224 {
		bits = append(bits, 224)
	} else if bit == 256 {
		bits = append(bits, 256)
	} else if bit == 384 {
		bits = append(bits, 384)
	} else if bit == 512 {
		bits = append(bits, 512)
	} else if bit == 512224 {
		bits = append(bits, 512224)
	} else if bit == 512256 {
		bits = append(bits, 512256)
	} else {
		fmt.Println("bit not supported, (1, 224, 256, 384, 512, 512224, 512256)")
		return
	}

	input, e := c.Flags().GetString("input")
	if e == nil && len(input) > 0 {
		fmt.Printf("%9s : %s\n", "file", input)
	} else {
		fmt.Printf("%9s : %s\n", "source", strings.Replace(string(content), "\n", "", -1))
	}

	for _, b := range bits {
		var name string
		if b == 1 {
			name = "sha1"
			handle = sha1.New()
		} else if b == 224 {
			name = "sha224"
			handle = sha256.New224()
		} else if b == 256 {
			name = "sha256"
			handle = sha256.New()
		} else if b == 384 {
			name = "sha384"
			handle = sha512.New384()
		} else if b == 512 {
			name = "sha512"
			handle = sha512.New()
		} else if b == 512224 {
			name = "sha512224"
			handle = sha512.New512_224()
		} else if b == 512256 {
			name = "sha512256"
			handle = sha512.New512_256()
		}
		handle.Write(content)
		value := hex.EncodeToString(handle.Sum(nil))
		fmt.Printf("%9s : %s\n", name, value)
		fmt.Printf("%9s : %s\n", "", strings.ToUpper(value))
	}

}
func init() {
	sha1Cmd.Flags().StringP("input", "i", "", "输入待验证的文件")
	sha1Cmd.Flags().IntP("bit", "b", 1, "sha bit: 1, 224, 256, 384, 512, 512224, 512256")
}

func NewShaCmd() *cobra.Command {
	return sha1Cmd
}
