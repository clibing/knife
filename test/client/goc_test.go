package client

import (
	"fmt"
	"testing"

	"github.com/clibing/knife/cmd/client"
)

func TestGoc(t *testing.T) {
	fmt.Printf("|%6d|%6d|\n", 12, 1234567)
}

func TestCreateDir(t *testing.T) {
	urls := []string{
		"https://github.com/clibing/knife.git",
		"https://gitea.linuxcrypt.cn/macOS/hysteria-app.git",
	}

	for _, url := range urls {
		value, _ := client.CreateDir(url)
		fmt.Println(value)
	}
}
