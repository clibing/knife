package server

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	port int
	path string
)

// staticCmd represents the static command
var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "文件服务。",
	Long: `启用本地静态资源服务:

新装系统后，安装所需软件的时候，每次都需要移动硬盘、U盘或者scp等拷贝资源到目标机器。
一般情况都有一台闲置的电脑, 被安装的电脑在安装机器的期间, 可以使用闲置的机器可以去官网现在所需最新的软件安装包。`,
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")

		templateData := make(map[string]bool)
		templateData["token"] = len(token) > 0

		http.HandleFunc("/upload", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					t, e := template.New("upload").Parse(`
<html><head><title>Upload file</title></head><body>
<form enctype="multipart/form-data" action="upload" method="post">
	{{ if .token }} 
	<input type="input" name="token" /> 上传凭证 <br/>
	{{ end }}
	<input type="name" name="path" /> 存储目录 <br/>
	<input type="file" name="file" /> 
	<input type="submit" value="upload" />
</form></body></html>`)
					if e != nil {
						fmt.Println(e)
						return
					}
					t.Execute(w, templateData)
					return
				}

				inputToken := r.FormValue("token")
				if len(token) > 0 && token != inputToken {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("凭证为空"))
					return
				}

				// parse and validate source and post parameters
				source, m, err := r.FormFile("file")
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("获取上传文件异常"))
					return
				}
				defer source.Close()

				// 文件名字
				storagePath := r.FormValue("path")
				storagePath = strings.TrimPrefix(storagePath, "/")
				name := filepath.Join(path, storagePath, m.Filename)
				parent := filepath.Dir(name)
				if _, e := os.Stat(parent); e != nil {
					if os.IsNotExist(e) {
						e = os.MkdirAll(parent, os.ModePerm)
						if e != nil {
							fmt.Println(e)
						}
					}
				}

				n, e := os.Create(name)
				if e != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("创建文件异常"))
				}
				defer n.Close()

				data, e := io.ReadAll(source)
				if e != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("读取上传文件异常"))
					return
				}

				if _, e := n.Write(data); e != nil {
					if e != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("写入文件异常"))
						return
					}
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("<html><title>上传成功</title><body>上传成功<br/><img src='%s/%s'></body></html>", storagePath, m.Filename)))
			}),
		)
		// 静态资源的目录
		fs := http.FileServer(http.Dir(path))
		// http 处理器
		http.Handle("/", http.StripPrefix("/", fs))
		// 获取端口
		value := fmt.Sprintf(":%d", port)
		// 建立监听
		listener, err := net.Listen("tcp", value)
		if err != nil {
			fmt.Println("建立监听异常, ", err.Error())
			return
		}

		fmt.Printf("服务启动中: http://0.0.0.0:%d \n", listener.Addr().(*net.TCPAddr).Port)
		err = http.Serve(listener, nil)
		if err != nil {
			fmt.Println("服务启动失败，请检查, ", err.Error())
			return
		}
	},
}

func init() {
	staticCmd.Flags().StringVarP(&path, "path", "p", "", "静态资源目录, 默认为当前目录")
	staticCmd.Flags().IntVarP(&port, "port", "", 0, "端口, 默认会随机")
	staticCmd.Flags().StringP("token", "t", "", "上传开启凭证, 当为空时，不启用")
}

func NewFileServer() *cobra.Command {
	return staticCmd
}
