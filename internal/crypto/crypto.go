// Package crypto 提供无 cobra 依赖的对称加密纯函数层：
// 函数不做 IO、不读全局状态；IV 与编码由调用方决定。
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// BlockSize aes 块大小（字节）
const BlockSize = aes.BlockSize

// ---------- PKCS5 ----------

// PKCS5Padding 补全码：数据长度不对齐时填充 n 个值为 n 的字节
func PKCS5Padding(content []byte, blockSize int) []byte {
	padding := blockSize - len(content)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(content, padtext...)
}

// PKCS5UnPadding 去除补全码；空数据或非法 padding 返回错误
func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("crypto: empty data")
	}
	unpadding := int(origData[length-1])
	if unpadding == 0 || unpadding > length {
		return nil, fmt.Errorf("crypto: invalid padding byte %d", unpadding)
	}
	return origData[:length-unpadding], nil
}

// ---------- CBC ----------

// CBCEncrypt CBC 加密，iv 长度必须等于 BlockSize
func CBCEncrypt(plain, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != BlockSize {
		return nil, fmt.Errorf("crypto: iv length %d, want %d", len(iv), BlockSize)
	}
	data := PKCS5Padding(plain, BlockSize)
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	mode.CryptBlocks(encrypted, data)
	return encrypted, nil
}

// CBCDecrypt CBC 解密，密文长度必须是块大小的正整数倍
func CBCDecrypt(encrypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != BlockSize {
		return nil, fmt.Errorf("crypto: iv length %d, want %d", len(iv), BlockSize)
	}
	if len(encrypted) == 0 || len(encrypted)%BlockSize != 0 {
		return nil, fmt.Errorf("crypto: ciphertext length %d is not a positive multiple of %d", len(encrypted), BlockSize)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)
	return PKCS5UnPadding(decrypted)
}

// ---------- ECB ----------

// GenerateKey 将任意长度 key 归一化为 16 字节（ECB 使用）：
// 不足补 0，超出部分折叠 XOR 进前 16 字节
func GenerateKey(key []byte) []byte {
	genKey := make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// ECBEncrypt ECB 加密（块级独立加密，无 IV）
func ECBEncrypt(origData []byte, key []byte) []byte {
	block, _ := aes.NewCipher(GenerateKey(key))
	length := (len(origData) + BlockSize) / BlockSize
	plain := make([]byte, length*BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	for bs, be := 0, BlockSize; bs < len(plain); bs, be = bs+BlockSize, be+BlockSize {
		block.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// ECBDecrypt ECB 解密
func ECBDecrypt(encrypted []byte, key []byte) []byte {
	block, _ := aes.NewCipher(GenerateKey(key))
	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, BlockSize; bs < len(encrypted); bs, be = bs+BlockSize, be+BlockSize {
		block.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

// ---------- CFB ----------

// CFBEncrypt CFB 加密：内部生成随机 IV，返回 IV||密文
func CFBEncrypt(plain, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, BlockSize+len(plain))
	iv := encrypted[:BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[BlockSize:], plain)
	return encrypted, nil
}

// CFBDecrypt CFB 解密：入参为 IV||密文
func CFBDecrypt(encrypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < BlockSize {
		return nil, errors.New("crypto: ciphertext too short")
	}
	iv := encrypted[:BlockSize]
	body := encrypted[BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(body, body)
	return body, nil
}

// ---------- 编码识别 ----------

// IsHex 判断字节序列是否全为 hex 字符（空序列视为 hex）
func IsHex(value []byte) bool {
	for _, v := range value {
		if (v >= '0' && v <= '9') || (v >= 'A' && v <= 'F') || (v >= 'a' && v <= 'f') {
			continue
		}
		return false
	}
	return true
}
