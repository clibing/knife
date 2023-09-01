package sign

/**
 * ZeroPadding，数据长度不对齐时使用0填充，否则不填充
 * PKCS7Padding，假设数据长度需要填充n(n>0)个字节才对齐，那么填充n个字节，每个字节都是n;如果数据本身就已经对齐了，则填充一块长度为块大小的数据，每个字节都是块大小
 * PKCS5Padding，PKCS7Padding的子集，块大小固定为8字节。
 * 两者的区别在于PKCS5Padding是限制块大小的PKCS7Padding.
 */

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

	if !skip_console {
		fmt.Printf("%14s : %s\n", "mode", mode)
	}
	if mode == "CBC" {
		var err error
		if !decrypt {
			result, err = cbcEncrypt(c, content, []byte(key))
		} else {
			result, err = cbcDecrypt(c, content, []byte(key))
		}
		if err != nil {
			fmt.Printf("%14s : %s\n", "err", err.Error())
			return
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

func iv(c *cobra.Command, key string, blockSize int) (iv []byte, err error) {
	input, _ := c.Flags().GetString("iv")
	if len(input) == len(key) && len(input) == blockSize {
		iv = []byte(input)
		return
	}

	if len(input) < blockSize {
		err = fmt.Errorf("偏移量小于块大小, 当前块大小: %d", blockSize)
		return
	} else if len(input) > blockSize {
		fmt.Printf("偏移量过长，按照块大小(%d)自动截取\n", blockSize)
		iv = []byte(input)[:blockSize]
	}
	err = fmt.Errorf("未知错误: %s", input)
	return
}
func cbcEncrypt(c *cobra.Command, content []byte, key []byte) (value []byte, err error) {
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()           // 获取秘钥块的长度
	data := pkcs5Padding(content, blockSize) // 补全码
	iv, err := iv(c, string(key), blockSize)
	if err != nil {
		return
	}
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted := make([]byte, len(data))           // 创建数组
	blockMode.CryptBlocks(encrypted, data)         // 加密
	value = encrypted
	return
}

func cbcDecrypt(c *cobra.Command, encrypted []byte, key []byte) (value []byte, err error) {
	block, _ := aes.NewCipher(key) // 分组秘钥
	blockSize := block.BlockSize() // 获取秘钥块的长度
	iv, err := iv(c, string(key), blockSize)
	if err != nil {
		return
	}
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted := make([]byte, len(encrypted))      // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)    // 解密
	decrypted = pkcs5UnPadding(decrypted)          // 去除补全码
	value = decrypted
	return
}

func pkcs5Padding(content []byte, blockSize int) []byte {
	padding := blockSize - len(content)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(content, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// =========================================
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
	aesCmd.Flags().String("iv", "", "当model为CBC时生效，偏移量可用，当为空时使用密钥")
	aesCmd.Flags().StringP("mode", "m", "CBC", "模式: CBC, ECB, CFB")
	aesCmd.Flags().BoolP("decrypt", "d", false, "编码方式，默认encrypt")
}

func NewAesCmd() *cobra.Command {
	return aesCmd
}
