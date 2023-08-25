package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func TestUtils(t *testing.T) {
	// flag := "http-debug"
	// clientCmd.Flags().Bool(flag, false, "")
	// messageAutoEntry(clientCmd, "it's error", "Not Found")
	// messageAutoEntry(clientCmd, "it's error: %s, %s", "Not Found, %s, %s", "1", "root")

	// flag = "websocket-debug"
	// clientCmd.Flags().Bool("websocket-debug", true, "")
	// messageAutoEntry(clientCmd, "it's error", "Not Found")
	// messageAutoEntry(clientCmd, "it's error: %s, %s", "Not Found: %s, %s", "1", "root")
}

func TestJson(t *testing.T) {
	v := make(map[string]interface{})
	v["admin"] = "admin"
	v["password"] = "123456"

	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}

func TestUrl(t *testing.T) {
	req := "www.baidu.com?v=100"
	u, _ := url.Parse(req)
	fmt.Println(u.Path)

}
