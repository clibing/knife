package test

// import (
// 	"bytes"
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"encoding/base64"
// 	"fmt"
// 	"testing"
// )

// func TestAesMehotd(t *testing.T) {
// 	v, _ := AesEncryptCBC([]byte("clibing"), []byte("1234567890abcedf"))
// 	fmt.Println(len(v))
// }

// // 补齐算法
// func PKCSPadding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(ciphertext, padtext...)
// }

// func PKCSUnPadding(origData []byte) []byte {
// 	length := len(origData)
// 	unpadding := int(origData[length-1])
// 	return origData[:(length - unpadding)]
// }

// // 注意：本示例中轮秘钥和初始化向量使用相同的字节数据，真实场景不推荐
// func AesEncryptCBC(origData, key []byte) ([]byte, error) {
// 	// AES算法的轮秘钥 key
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		fmt.Errorf("encrypt error %v", err)
// 		return nil, err
// 	}

// 	blockSize := block.BlockSize()
// 	// origData = PKCSPadding(origData, blockSize)
// 	// key[:blockSize]初始化向量
// 	origData = PKCSPadding(origData, blockSize)
// 	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
// 	crypted := make([]byte, len(origData))
// 	blockMode.CryptBlocks(crypted, origData)
// 	fmt.Printf(" encrypt base64 string  %s\n", base64.StdEncoding.EncodeToString(crypted))

// 	return crypted, nil
// }

// func AesDecryptCBC(crypted, key []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		fmt.Errorf(" encrypt error %v", err)
// 		return nil, err
// 	}

// 	blockSize := block.BlockSize()
// 	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
// 	origData := make([]byte, len(crypted))
// 	blockMode.CryptBlocks(origData, crypted)
// 	origData = PKCSUnPadding(origData)
// 	return origData, nil
// }
