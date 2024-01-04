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

	"github.com/clibing/go-common/pkg/util"
	"github.com/go-basic/uuid"
	"github.com/spf13/cobra"
)

// const maxUploadSize = 2 * 1024 * 1024 // 2 mb

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

		templateData := make(map[string]bool)
		token, _ := cmd.Flags().GetString("token")
		templateData["token"] = len(token) > 0

		// code
		code := make(map[string]string)
		http.HandleFunc("/verfiy", func(w http.ResponseWriter, r *http.Request) {
			param := r.URL.Query()
			input := param.Get("token")

			w.Header().Set("Content-Type", "application/json")
			// token 是否开启，
			if len(token) == 0 {
				// 没有开启
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{'code': 201, 'message': '暂未开启token认证'}"))
				return
			}

			if len(token) > 0 && token != input {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("{'code': 401, 'message': '凭证错误'}"))
				return
			}

			// 去除 - 短接线
			key := strings.ReplaceAll(uuid.New(), "-", "")
			// 存储cache
			code[key] = ""
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("{'code': 200, 'message': '认证成功', 'data': '%s'}", key)))
		})

		http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				t, e := template.New("upload").Parse(`
<html><head><title>Upload file</title></head><body>
<form enctype="multipart/form-data" action="upload" method="post">
{{ if .token }} 
<input type="input" name="token" /> 上传凭证 <br/>
{{ end }}
<input type="name" name="path" /> 存储目录 <br/>
<input type="button" value="+"/><br/>
<ol>
	<li>
		<input type="file" name="file" /> <br/>
	</li>
</ol>
<input type="submit" value="upload" />
</form></body></html>`)
				if e != nil {
					fmt.Println(e)
					return
				}
				t.Execute(w, templateData)
				return
			}

			// 32<<20  ==> 32MB
			// 32<<21  ==> 64MB
			// 32<<25  ==> 1024MB
			maxMemory, _ := cmd.Flags().GetString("maxMemory")
			err := r.ParseMultipartForm(util.ReverseByteFormat(maxMemory)) // 设置最大上传文件大小为32MB
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("允许最大上传文件为: %s", maxMemory)))
				return
			}
			// 文件名字
			storagePath := r.FormValue("path")
			storagePath = strings.TrimPrefix(storagePath, "/")

			key := r.FormValue("token")
			if _, ok := code[key]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("验证中间token不存在"))
				return
			}
			// 删除 key
			delete(code, key)

			multiFiles := r.MultipartForm.File["file"]
			for _, m := range multiFiles {

				// parse and validate source and post parameters
				source, err := m.Open()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("获取上传文件异常"))
					return
				}
				defer source.Close()

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
			}
			http.Redirect(w, r, "/", http.StatusOK)
		})
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
	staticCmd.Flags().StringP("maxMemory", "m", "32M", "设置内存大小")
}

func NewFileServer() *cobra.Command {
	return staticCmd
}
