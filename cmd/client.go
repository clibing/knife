/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"moul.io/http2curl"
)

const (
	CMD_CLIENT_HTTP_NAME         string = "http"
	CMD_CLIENT_HTTP_METHOD       string = "method"
	CMD_CLIENT_HTTP_URL          string = "url"
	CMD_CLIENT_HTTP_SHOW_CURL    string = "show-curl"
	CMD_CLIENT_HTTP_CONTENT_TYPE string = "Content-Type"
	CMD_CLIENT_HTTP_PATH         string = "path"
)

// clientCmd represents the client command
var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "client",
		Long:  `多应用客户端`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	httpCmd = &cobra.Command{
		Use:   CMD_CLIENT_HTTP_NAME,
		Short: CMD_CLIENT_HTTP_NAME,
		Long:  `http client`,
		Run: func(c *cobra.Command, args []string) {
			debug := NewDebug(c)

			var err error
			url, err := c.Flags().GetString(CMD_CLIENT_HTTP_URL)
			if err != nil {
				debug.ShowSame("method is required, err: %s", err.Error())
				return
			}

			if len(url) == 0 {
				debug.ShowSame("method is required, url is empty")
				return
			}
			url = formatURL(url)

			method, _ := c.Flags().GetString(CMD_CLIENT_HTTP_METHOD)
			if len(method) == 0 {
				debug.ShowSame("method is required")
				return
			}

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Timeout: 15 * time.Second, Transport: tr}

			var request *http.Request

			// 获取  request header
			headers := ParseHeader(c)
			// 格式 请求数据
			body := ParseData(c, headers, debug)

			request, err = http.NewRequest(method, url, body)
			for k, v := range headers {
				request.Header.Add(k, v)
			}

			if err != nil {
				debug.ShowSame("create request failed. err: %s", err.Error())
				return
			}

			if showCurl, err := c.Flags().GetBool(CMD_CLIENT_HTTP_SHOW_CURL); err == nil && showCurl {
				v, _ := http2curl.GetCurlCommand(request)
				debug.ShowSame("%s", v)
			}

			response, err := client.Do(request)
			if err != nil {
				debug.ShowSame("send request failed. err: %s", err.Error())
				return
			}
			defer response.Body.Close()

			data, _ := io.ReadAll(response.Body)
			path, enable, _ := IsDownload(c)
			if enable {
				os.WriteFile(path, data, 0644)
			} else {
				content_type := response.Header.Get(CMD_CLIENT_HTTP_CONTENT_TYPE)
				if isText(content_type) || isJson(content_type) {
					debug.ShowSame("%s", response.Status)
					debug.ShowSame("%s", string(data))
				}
			}

		},
	}
	websocketCmd = &cobra.Command{
		Use:     "websocket",
		Short:   "websocket",
		Aliases: []string{"ws"},
		Long:    `websocket client`,
		Run: func(ws *cobra.Command, args []string) {
			fmt.Println("websocket client called")
		},
	}
)

func IsDownload(c *cobra.Command) (path string, enable bool, err error) {
	path, err = c.Flags().GetString(CMD_CLIENT_HTTP_PATH)
	if err != nil {
		return
	}
	enable = len(path) > 0
	return
}

// 格式化 data 数据
func ParseData(c *cobra.Command, headers map[string]string, debug *Debug) (body io.Reader) {
	data, err := c.Flags().GetStringSlice("data")
	if err == nil {
		if len(data) > 0 {
			if contentType, ok := headers[CMD_CLIENT_HTTP_CONTENT_TYPE]; ok && isJson(contentType) {
				p := DataToJson(data)
				debug.ShowSame("%s", p)
				body = strings.NewReader(p)
			} else {
				p := strings.Join(data, "&")
				debug.ShowSame("%s", p)
				body = strings.NewReader(p)
			}
		}
	}
	return
}

func DataToJson(data []string) (result string) {
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

func ParseHeader(c *cobra.Command) (result map[string]string) {
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
	return false
}

func isJson(contentType string) bool {
	if len(contentType) == 0 {
		return false
	}

	if strings.HasPrefix(contentType, "application/json") {
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

func init() {
	httpCmd.Flags().StringP(CMD_CLIENT_HTTP_PATH, "p", "", "save http response body to file path")
	httpCmd.Flags().Bool(fmt.Sprintf("%s-debug", httpCmd.Use), false, "set debug, default: false")
	httpCmd.Flags().Bool(CMD_CLIENT_HTTP_SHOW_CURL, false, "show curl execute bash, default: false")
	httpCmd.Flags().StringP(CMD_CLIENT_HTTP_URL, "u", "", "request url")
	httpCmd.Flags().StringP(CMD_CLIENT_HTTP_METHOD, "m", "GET", "request method [GET|POST|...], default GET")
	httpCmd.Flags().StringSliceP("header", "H", []string{}, "header, 格式: \"-H 'key: value'\"")
	httpCmd.Flags().StringSliceP("data", "d", []string{}, "data key=value")

	// client -> http request
	clientCmd.AddCommand(httpCmd)

	websocketCmd.Flags().Bool("websocket-debug", false, "debug, default false")
	// client -> websocket
	clientCmd.AddCommand(websocketCmd)

	rootCmd.AddCommand(clientCmd)

	// clientCmd.Flags().StringVarP(&clientType, "client-type", "c", "", "选择客户端类型： http, socks5, websocket")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
