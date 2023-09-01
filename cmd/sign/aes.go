package sign

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var aesCmd = &cobra.Command{
	Use:   "aes",
	Short: `aes加密解密`,
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
	if show_file {
		fmt.Printf("%14s : %s\n", "source", output)
	} else {
		fmt.Printf("%14s : %s\n", "source", strings.Replace(string(content), "\n", "", -1))
	}

	if decrypt {
		var e error
		if isHex(content) {
			content, e = hex.DecodeString(string(content))
		} else {
			content, e = base64.StdEncoding.DecodeString(string(content))
		}
		if e != nil {
			fmt.Println("待解密内容异常", e.Error())
			return
		}
	}

	fmt.Printf("%14s : %s\n", "mode", mode)
	if mode == "CBC" {
		if !decrypt {
			result = cbcEncrypt(content, []byte(key))
		} else {
			result = cbcDecrypt(content, []byte(key))
		}
	} else if mode == "ECB" {
		if !decrypt {
			result = ecbEncrypt(content, []byte(key))
		} else {
			result = ecbDecrypt(content, []byte(key))
		}
	} else if mode == "CFB" {
		if !decrypt {
			result = cfbEncrypt(content, []byte(key))
		} else {
			result = cfbDecrypt(content, []byte(key))
		}
	} else {
		fmt.Println("Mode not supported: [CBC, ECB, CFB]", mode)
		return
	}

	fmt.Printf("%14s : %s\n", "key", key)
	if decrypt {
		fmt.Printf("%14s : %s\n", "value", string(result))
	} else {
		value := hex.EncodeToString(result)
		value2 := base64.StdEncoding.EncodeToString(result)
		fmt.Printf("%14s : %s\n", "value(hex)", value)
		fmt.Printf("%14s : %s\n", "value(base64)", value2)

	}

	verify, _ := c.Flags().GetString("verify")
	if len(verify) > 0 {
		fmt.Printf("%14s : %s\n", "verify", verify)
	}
}

func cbcEncrypt(content []byte, key []byte) (value []byte) {
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                // 获取秘钥块的长度
	padding := blockSize - len(content)%blockSize // 补全码
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	data := append(content, padtext...)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted := make([]byte, len(data))                        // 创建数组
	blockMode.CryptBlocks(encrypted, data)                      // 加密
	return encrypted
}

func cbcDecrypt(encrypted []byte, key []byte) (value []byte) {
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted := make([]byte, len(encrypted))                   // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密

	length := len(decrypted)
	unpadding := int(decrypted[length-1])
	return decrypted[:(length - unpadding)]
}

func ecbEncrypt(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func ecbDecrypt(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

func cfbEncrypt(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}
func cfbDecrypt(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

func isHex(value []byte) bool {
	for _, v := range value {
		if v >= 48 && v <= 57 || v >= 65 && v <= 70 || v >= 97 && v <= 102 {
			continue
		} else {
			return false
		}
	}
	return true
}

func init() {
	aesCmd.Flags().StringP("key", "k", "", "加密的密钥")
	aesCmd.Flags().StringP("input", "i", "", "输入文件")
	aesCmd.Flags().StringP("output", "o", "", "输出的文件")
	aesCmd.Flags().StringP("verify", "v", "", "待验证的信息")
	aesCmd.Flags().StringP("mode", "m", "CBC", "模式: CBC, ECB, CFB")
	aesCmd.Flags().BoolP("decrypt", "d", false, "编码方式，默认encrypt")
}

func NewAesCmd() *cobra.Command {
	return aesCmd
}
