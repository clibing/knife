package sign

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/spf13/cobra"
)

var bits int

// rsaCmd represents the rsa command
var rsaCmd = &cobra.Command{
	Use:   "rsa",
	Short: "公钥私钥",
	Long: `快速生成公钥、私钥:

1. 生成密钥长度为1024
knife rsa -b 1024
.`,
	Run: func(_ *cobra.Command, _ []string) {
		genRsaKey(bits)
	},
}

func init() {
	rsaCmd.Flags().IntVarP(&bits, "bits", "b", 2048, "密钥长度，默认为1024位")
}

func NewRsaCmd() *cobra.Command {
	return rsaCmd
}

// RSA公钥私钥产生
func genRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	privatePKCS8Key, _ := x509.MarshalPKCS8PrivateKey(privateKey)

	blockPKCS8 := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privatePKCS8Key,
	}
	filePKCS8, err := os.Create("private.key")
	if err != nil {
		return err
	}
	err = pem.Encode(filePKCS8, blockPKCS8)
	if err != nil {
		return err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}
