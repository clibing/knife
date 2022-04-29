/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

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
	Run: func(cmd *cobra.Command, args []string) {
		genRsaKey(bits)
	},
}

func init() {
	rootCmd.AddCommand(rsaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rsaCmd.PersistentFlags().String("foo", "", "A help for foo")
	rsaCmd.Flags().IntVarP(&bits, "bits", "b", 1024, "密钥长度，默认为1024位")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rsaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//RSA公钥私钥产生
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

	privatePKCS8Key, err := x509.MarshalPKCS8PrivateKey(privateKey)

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
