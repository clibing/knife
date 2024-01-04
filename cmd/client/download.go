package client

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"runtime"
	"strings"
	"time"

	"github.com/clibing/knife/pkg/download"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "knife client download -u \"文件url\"",
	Run: func(c *cobra.Command, args []string) {
		url, err := c.Flags().GetString("url")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(url) == 0 {
			fmt.Println("url为空")
			return
		}

		output, err := c.Flags().GetString("output")
		dir, err := c.Flags().GetString("dir")
		headers, err := c.Flags().GetStringSlice("headers")
		ck, err := c.Flags().GetString("cookie-key")
		cv, err := c.Flags().GetString("cookie-value")

		sm, err := c.Flags().GetString("signature-method")
		sv, err := c.Flags().GetString("signature-value")
		if len(sv) > 0 {
			if sm != "md5" && sm != "sha1" && sm != "sha256" && sm != "sha512" {
				fmt.Println("当前签名类型错误: ", sm, " 目前支持md5, sha1, sha256, sha512")
				return
			}
		}
		var signatureMethod func([]byte) string
		if len(sv) > 0 {
			signatureMethod = func(b []byte) string {
				var m hash.Hash
				switch sm {
				case "md5":
					m = md5.New()
				case "sha1":
					m = sha1.New()
				case "sha256":
					m = sha256.New()
				case "sha512":
					m = sha512.New()
				}
				m.Write(b)
				return hex.EncodeToString(m.Sum(nil))
			}
		}
		task, err := c.Flags().GetInt("task")

		h := make(map[string]string)
		if len(headers) > 0 {
			for _, v := range headers {
				c := strings.Split(v, "=")
				key := strings.Trim(c[0], "")
				value := strings.Trim(c[1], "")
				h[key] = value
			}
		}

		start := time.Now()
		fmt.Println("开启: ", start.Format("2006-01-02 15:04:05"))
		err = download.NewFileDownloader(url, output, dir, sv, task, signatureMethod, h, ck, cv).Run()
		if err != nil {
			fmt.Println("下载任务失败: ", err.Error())
		} else {
			end := time.Now()
			d := end.UnixMicro() - start.UnixMicro()
			// util.
			fmt.Printf("耗时: %s (%s)\n", end.Format("2006-01-02 15:04:05"), format(d))
		}
	},
}

func init() {
	downloadCmd.Flags().StringP("url", "u", "", "下载文件的url")
	downloadCmd.Flags().StringP("output", "o", "", "可选，输出文件的名字，默认读取header头")
	downloadCmd.Flags().StringP("dir", "d", "", "可选，指定文件夹, 默认为当前目录")
	downloadCmd.Flags().IntP("task", "t", runtime.NumCPU(), "可选，默认分片下载的任务数, 默认为当前核心数")
	downloadCmd.Flags().StringP("signature-value", "v", "", "可选，待验证文件的值，为空时不校验")
	downloadCmd.Flags().StringP("signature-method", "m", "", "可选，文件签名方法, 可选 md5 sha1 sha256 sha512")
	downloadCmd.Flags().StringSliceP("headers", "H", []string{}, "可选，下载文件携带请求头, 例如: \"-H Token=abcdef\"")
	downloadCmd.Flags().StringP("cookie-key", "k", "", "可选，下载时携带的cookie的key")
	downloadCmd.Flags().StringP("cookie-value", "c", "", "可选，下载时携带的cookie的value")
}

func NewDownloadCmd() *cobra.Command {
	return downloadCmd
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
