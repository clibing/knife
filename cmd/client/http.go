package client

import (
	"encoding/json"
	"net/url"
	"strings"

	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/clibing/knife/cmd/debug"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
	"moul.io/http2curl"
)

const (
	CMD_CLIENT_HTTP_NAME         string = "http"
	CMD_CLIENT_HTTP_METHOD       string = "method"
	CMD_CLIENT_HTTP_SHOW_CURL    string = "show-curl"
	CMD_CLIENT_HTTP_INPUT_TYPE   string = "input-type"
	CMD_CLIENT_HTTP_CONTENT_TYPE string = "Content-Type"
	CMD_CLIENT_HTTP_OUTPUT       string = "output"
)

var (
	httpCmd = &cobra.Command{
		Use:   CMD_CLIENT_HTTP_NAME,
		Short: CMD_CLIENT_HTTP_NAME,
		Long:  `http client`,
		Example: `1.常用
knife client http https://tool.linuxcrypt.cn/checkRemoteIp

2.将Response.body保存到文件中/tmp/result.data
knife client http https://tool.linuxcrypt.cn/checkRemoteIp -o /tmp/result.data

3.请求JSON数据
knife client http https://tool.linuxcrypt.cn/checkRemoteIp -m POST -d '{}'  -H 'Content-Type: application/json; charset=utf-8'

4.非常规提交数据, 格式为key=value, 当请求的header为json，需要设置--form-to-json，会自动转换。
client http https://tool.linuxcrypt.cn/checkRemoteIp -m POST -d 'name=admin' -d 'root=123456'  -H 'Content-Type: application/json; charset=utf-8' --form-to-json

5.当请求时，需要转换curl格式
knife client http https://tool.linuxcrypt.cn/checkRemoteIp --show-curl`,
		Run: func(c *cobra.Command, args []string) {
			debug := debug.NewDebug(c)

			if len(args) == 0 {
				c.Help()
				return
			}

			for _, url := range args {
				var err error
				if len(url) == 0 {
					debug.ShowSame("url is empty")
					continue
				}
				url = formatURL(url)

				method, _ := c.Flags().GetString(CMD_CLIENT_HTTP_METHOD)
				if len(method) == 0 {
					debug.ShowSame("method is required")
					continue
				}

				tr := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
				client := &http.Client{Timeout: 15 * time.Second, Transport: tr}

				var request *http.Request

				// 获取  request header
				headers := parseHeader(c)
				// 格式 请求数据
				body, err := parseRequestData(c, headers, debug)
				if err != nil {
					debug.ShowSame("%s", err.Error())
					return
				}

				request, err = http.NewRequest(method, url, body)
				for k, v := range headers {
					request.Header.Add(k, v)
				}

				if err != nil {
					debug.ShowSame("create request failed. err: %s", err.Error())
					continue
				}

				if showCurl, err := c.Flags().GetBool(CMD_CLIENT_HTTP_SHOW_CURL); err == nil && showCurl {
					v, _ := http2curl.GetCurlCommand(request)
					debug.ShowSame("%s", v)
				}

				response, err := client.Do(request)
				if err != nil {
					debug.ShowSame("send request failed. err: %s", err.Error())
					continue
				}
				defer response.Body.Close()

				data, _ := io.ReadAll(response.Body)
				path, enable, _ := isDownload(c)
				if enable {
					os.WriteFile(path, data, 0644)
				} else {
					content_type := response.Header.Get(CMD_CLIENT_HTTP_CONTENT_TYPE)
					if isText(content_type) || isJson(content_type) {
						debug.ShowSame("%s\n", response.Status)
						if isJson(content_type) {
							beautify, e := c.Flags().GetBool("beautify")
							if e == nil && beautify {
								option := &pretty.Options{Width: 80, Prefix: "", Indent: "  ", SortKeys: false}
								value := pretty.PrettyOptions([]byte(data), option)
								debug.ShowSame("%s", string(value))
							} else {
								debug.ShowSame("%s", string(data))
							}
						} else {
							debug.ShowSame("%s", string(data))
						}
					}
				}
			}
		},
	}
)

func isText(contentType string) bool {
	if len(contentType) == 0 {
		return false
	}
	if strings.HasPrefix(contentType, "text") {
		return true
	}
	if strings.Contains(contentType, "application/json") {
		return true
	}
	if strings.Contains(contentType, "text/json") {
		return true
	}
	return false
}

func isJson(contentType string) bool {
	if len(contentType) == 0 {
		return false
	}

	if strings.HasPrefix(contentType, "application/json") {
		return true
	}
	if strings.HasPrefix(contentType, "text/json") {
		return true
	}
	return false
}

func formatURL(request string) string {
	u, e := url.Parse(request)
	if e != nil {
		return request
	}

	scheme := u.Scheme
	path := u.Path
	query := u.RawQuery
	if len(scheme) == 0 {
		scheme = "http"
	}
	if len(path) == 0 {
		return scheme + "://" + u.Host
	}
	if len(query) == 0 {
		return scheme + "://" + u.Host + path
	}
	return scheme + "://" + u.Host + path + "?" + query
}

func dataToJson(data []string) (result string) {
	if len(data) == 0 {
		return
	}
	m := make(map[string]interface{})
	for _, v := range data {
		key, value, found := strings.Cut(v, "=")
		if found {
			m[key] = value
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	result = string(b)
	return
}

func parseHeader(c *cobra.Command) (result map[string]string) {
	result = make(map[string]string)
	// 请求 时 header 扩展
	headers, err := c.Flags().GetStringSlice("header")
	if err != nil || len(headers) == 0 {
		return
	}

	for _, header := range headers {
		key, value, found := strings.Cut(header, ": ")
		if found {
			result[key] = value
		}
	}
	return
}

func isDownload(c *cobra.Command) (path string, enable bool, err error) {
	path, err = c.Flags().GetString(CMD_CLIENT_HTTP_OUTPUT)
	if err != nil {
		return
	}
	enable = len(path) > 0
	return
}

// 格式化 data 数据
func parseRequestData(c *cobra.Command, headers map[string]string, debug *debug.Debug) (body io.Reader, err error) {
	data, err := c.Flags().GetStringSlice("data")
	if err == nil {
		if len(data) > 0 {
			if contentType, ok := headers[CMD_CLIENT_HTTP_CONTENT_TYPE]; ok && isJson(contentType) {
				var p string
				var b bool
				if b, err = c.Flags().GetBool("form-to-json"); b && err == nil {
					p = dataToJson(data)
				} else {
					size := len(data)
					if size > 1 {
						err = fmt.Errorf("按application/json提交,数据存在多段")
						return
					} else if size == 0 {
						return
					} else {
						p = data[0]
						if !json.Valid([]byte(p)) {
							err = fmt.Errorf("输入的data不是json格式")
							return
						}
					}
				}
				// debug.ShowSame("request body: %s", p)
				body = strings.NewReader(p)
			} else {
				p := strings.Join(data, "&")
				// debug.ShowSame("request body: %s", p)
				body = strings.NewReader(p)
			}
		}
	}
	return
}

func NewHttpClient() *cobra.Command {
	return httpCmd
}

func init() {
	httpCmd.Flags().StringSliceP("data", "d", []string{}, "data 数据，格式: \"-d username=admin\"")
	httpCmd.Flags().BoolP("form-to-json", "F", false, "输入form格式的数据自动转换为json提交，否则提交提交")
	httpCmd.Flags().BoolP("beautify", "b", true, "美化输出json, --beautify=false(忽略美化)")
	httpCmd.Flags().StringSliceP("header", "H", []string{}, "header, 格式: \"-H 'Content-Type: application/json; charset=utf-8'\"")
	httpCmd.Flags().BoolP(fmt.Sprintf("%s-debug", httpCmd.Use), "D", false, "是否启用debug, 默认: false")
	httpCmd.Flags().StringP(CMD_CLIENT_HTTP_METHOD, "m", "GET", "method [GET(获取资源)|HEAD(包头信息)|POST(增加资源)|PUT(更新-全字段)|PATCH(更新-目标字段)|DELETE(删除)|CONNECT|OPTIONS(获取支持的Method)]")
	httpCmd.Flags().StringP(CMD_CLIENT_HTTP_OUTPUT, "o", "", "将响应保存到指定文件")
	httpCmd.Flags().Bool(CMD_CLIENT_HTTP_SHOW_CURL, false, "是否展示curl命令行, 默认: false")
}
