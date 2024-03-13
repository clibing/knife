package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHttp2Client(t *testing.T) {
	ollama := "10.228.38.175:11434"

	data := map[string]string{
		"model":  "llama2-chinese",
		"prompt": "帮我创作一首诗",
	}

	body, _ := json.Marshal(data)

	tr := &http.Transport{
		MaxIdleConns: 100,
		// Dial: func(netw, addr string) (net.Conn, error) {
		// 	conn, err := net.DialTimeout(netw, addr, time.Second*120) //设置建立连接超时
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	err = conn.SetDeadline(time.Now().Add(time.Second * 120)) //设置发送接受数据超时
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	return conn, nil
		// },
		Dial: (&net.Dialer{
			Timeout:   120 * time.Second,
			KeepAlive: 120 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 120 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
	}

	req := &http.Request{
		Method: "POST",
		Header: http.Header{},
		URL: &url.URL{
			Scheme: "http",
			Host:   ollama,
			Path:   "/api/generate",
		},
		Body: io.NopCloser(bytes.NewReader(body)),
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(0),
	}

	response, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	if response.StatusCode == 500 {
		return
	}

	defer response.Body.Close()

	bufferedReader := bufio.NewReader(response.Body)

	buffer := make([]byte, 4*1024)

	var totalBytesReceived int

	type Response struct {
		Model              string    `json:"model"`
		CreatedAt          time.Time `json:"created_at"`
		Response           string    `json:"response"`
		Done               bool      `json:"done"`
		Context            []int     `json:"context"`
		PromptEvalDuration int64     `json:"prompt_eval_duration"`
		TotalDuration      int64     `json:"total_duration"`
		LoadDuration       int64     `json:"load_duration"`
		EvalCount          int       `json:"eval_count"`
		EvalDuration       int64     `json:"eval_duration"`
	}

	text := ""
	for {
		len, err := bufferedReader.Read(buffer)
		if len > 0 {
			totalBytesReceived += len
			value := buffer[0:len]
			fmt.Println(string(value))
			var response Response
			json.Unmarshal(value, &response)
			if !response.Done {
				text = text + response.Response

			}
		}

		if err != nil {
			if err == io.EOF {
				// Last chunk received
				// fmt.Println(err)
			}
			break
		}
	}
	fmt.Println(text)

}
