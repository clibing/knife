package download

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"runtime"
	"testing"
	"time"

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

	start := "2024-01-04 15:26:09"
	end := "2024-01-04 15:28:27"

	s, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	e, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	d := e.UnixMilli() - s.UnixMilli()
	fmt.Println(d)
	fmt.Println(format(d))
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
