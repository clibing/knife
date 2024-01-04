package download

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"runtime"
	"testing"

	"github.com/clibing/knife/pkg/download"
)

func TestDonwload(t *testing.T) {
	url := "https://discovery.linuxcrypt.cn/logo145x80.png"
	process := runtime.NumCPU()
	sign := "0170341a68ba2d86a951338b02a539b3"
	download.NewFileDownloader(url, "output.png", "/tmp", sign, process, func(data []byte) string {
		m := md5.New()
		m.Write(data)
		value := hex.EncodeToString(m.Sum(nil))
		fmt.Println("md5 value: ", value)
		return value
	}, nil, "", "").Run()
}

func TestFormat(*testing.T) {
	fmt.Println("800ms: ", format(800))
	fmt.Println("1500ms: ", format(1500))
	fmt.Println("180s: ", format(1000*60*3))
	fmt.Println("1.5h:  ", format(1000*60*90))
	fmt.Println("1.5d:  ", format(1000*60*90*24))
	fmt.Println("1.5m:  ", format(1000*60*90*24*30))
}

func format(micro int64) string {
	if micro < 1000 {
		return "不足1秒"
	} else if micro < 1000*60 {
		return fmt.Sprintf("%.2f 秒", float64(micro)/float64(1000))
	} else if micro < 1000*60*60 {
		return fmt.Sprintf("%.2f 分", float64(micro)/float64(1000*60))
	} else if micro < 1000*60*60*24 {
		return fmt.Sprintf("%.2f 时", float64(micro)/float64(1000*60*60))
	} else if micro < 1000*60*60*24*30 {
		return fmt.Sprintf("%.2f 天", float64(micro)/float64(1000*60*60*24))
	} else if micro < 1000*60*60*24*30*12 {
		return fmt.Sprintf("%.2f 月", float64(micro)/float64(1000*60*60*24*30))
	}
	return ""
}
