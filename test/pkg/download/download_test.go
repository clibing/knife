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
	}).Run()
}
