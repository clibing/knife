package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/clibing/go-common/pkg/util"
	"github.com/clibing/knife/internal/utils"
	"github.com/go-basic/uuid"
	"github.com/spf13/cobra"
)

var (
	port int
	path string
)

// 用户发现时携带的UA
const USER_AGENT = "clibing/knife"

type Result[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

func CreateErrMsg(code int, message string) (output []byte) {
	result := &Result[any]{
		Code:    code,
		Message: message,
	}
	output, _ = json.Marshal(result)
	return
}

func CreateSuccessMsg[T any](message string, data T) (output []byte) {
	result := &Result[T]{
		Code:    200,
		Data:    data,
		Message: message,
	}
	output, _ = json.Marshal(result)
	return
}

// staticCmd represents the static command
var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "文件服务。",
	Long: `启用本地静态资源服务:

新装系统后，安装所需软件的时候，每次都需要移动硬盘、U盘或者scp等拷贝资源到目标机器。
一般情况都有一台闲置的电脑, 被安装的电脑在安装机器的期间, 可以使用闲置的机器可以去官网现在所需最新的软件安装包。`,
	Run: func(cmd *cobra.Command, args []string) {

		found, _ := cmd.Flags().GetBool("found")
		if found {
			if port == 0 {
				fmt.Println("暂未设置静态服务端口.")
				return
			}
			foundStaticServer(port)
		}

		templateData := make(map[string]bool)
		token, _ := cmd.Flags().GetString("token")
		templateData["token"] = len(token) > 0

		// code
		code := make(map[string]string)
		http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
			param := r.URL.Query()
			input := param.Get("token")

			w.Header().Set("Content-Type", "application/json")
			// token 是否开启，
			if len(token) == 0 {
				// 没有开启
				w.WriteHeader(http.StatusOK)
				w.Write(CreateErrMsg(201, "暂未开启token认证"))
				return
			}

			if len(token) > 0 && token != input {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(CreateErrMsg(401, "凭证错误"))
				return
			}

			// 去除 - 短接线
			key := strings.ReplaceAll(uuid.New(), "-", "")
			// 存储cache
			code[key] = ""
			w.WriteHeader(http.StatusOK)
			w.Write(CreateSuccessMsg("认证成功", key))
		})

		http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				t, e := template.New("upload").Parse(HTML)
				if e != nil {
					fmt.Println(e)
					return
				}
				t.Execute(w, templateData)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			// 32<<20  ==> 32MB
			// 32<<21  ==> 64MB
			// 32<<25  ==> 1024MB
			maxMemory, _ := cmd.Flags().GetString("maxMemory")
			err := r.ParseMultipartForm(util.ReverseByteFormat(maxMemory)) // 设置最大上传文件大小为32MB
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(CreateErrMsg(400, fmt.Sprintf("允许最大上传文件为: %s", maxMemory)))
				return
			}
			// 文件名字
			storagePath := r.FormValue("path")
			storagePath = strings.TrimPrefix(storagePath, "/")

			key := r.FormValue("token")
			if _, ok := code[key]; !ok && len(token) > 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(CreateErrMsg(400, "验证中间token不存在"))
				return
			}
			// 删除 key
			delete(code, key)

			multiFiles := r.MultipartForm.File["file"]
			values := make([]string, 0)
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
				values = append(values, name)
			}

			ajax := r.FormValue("ajax")
			if ajax != "1" {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(CreateSuccessMsg("上传成功", values))
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

		// 获取 本地ip地址
		ip, _ := utils.GetLocalIp()
		port := listener.Addr().(*net.TCPAddr).Port
		// 发现服务
		go discovery(ip, port)

		fmt.Printf("服务启动中: http://%s:%d \n", ip, port)
		fmt.Printf("服务启动中: http://0.0.0.0:%d \n", port)
		err = http.Serve(listener, nil)
		if err != nil {
			fmt.Println("服务启动失败，请检查, ", err.Error())
			return
		}
	},
}

type Message struct {
	UserAgent string `json:"UserAgent"`
	Port      int    `json:"Port"`
}

// 注册器
func discovery(ip string, port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		var buf [512]byte
		n, remoteAddr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			continue
		}
		var msg Message
		json.Unmarshal(buf[:n], &msg)
		if msg.UserAgent == USER_AGENT {
			// 将本机ip+port发送给客户端
			SendIpAndPort(ip, remoteAddr.IP.String(), port, msg.Port)
		}
	}
}

func SendIpAndPort(localIp, remoteIp string, localPort, remotePort int) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", remoteIp, remotePort))
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(fmt.Sprintf("%s:%d", localIp, localPort)))
	if err != nil {
		return
	}
	fmt.Printf("新的客户端: %s:%d\n", remoteIp, remotePort)
}

func init() {
	staticCmd.Flags().StringVarP(&path, "path", "p", "", "静态资源目录, 默认为当前目录")
	staticCmd.Flags().IntVarP(&port, "port", "", 0, "端口, 默认会随机")
	staticCmd.Flags().StringP("token", "t", "", "上传开启凭证, 当为空时，不启用")
	staticCmd.Flags().StringP("maxMemory", "m", "", "设置内存大小, 默认不限制")
	staticCmd.Flags().BoolP("found", "f", false, "自动发现局域网内的静态服务器")
}

func NewFileServer() *cobra.Command {
	return staticCmd
}

func foundStaticServer(staticServerPort int) {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"), // 本地地址
		Port: 0,                      // 让操作系统随机选择端口
	}
	// 在UDP地址上建立UDP监听,得到连接
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	go func(nuc *net.UDPConn) {
		// 建立缓冲区
		stash := make(map[string]string, 0)
		buffer := make([]byte, 1024)
		for {
			//从连接中读取内容,丢入缓冲区
			i, udpAddr, e := nuc.ReadFromUDP(buffer)
			// 第一个是字节长度,第二个是udp的地址
			if e != nil {
				continue
			}
			key := string(buffer[:i])
			if _, ok := stash[key]; !ok {
				value := fmt.Sprintf("http://%s", key)
				fmt.Println(value)
				stash[key] = key
			}
			// 向客户端返回消息
			nuc.WriteToUDP([]byte("\n"), udpAddr)
		}
	}(conn)

	listenPort := conn.LocalAddr().(*net.UDPAddr).Port

	broadConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: staticServerPort,
	})
	if err != nil {
		fmt.Println("发送广播异常", err)
		return
	}
	defer broadConn.Close()

	msg := &Message{
		UserAgent: USER_AGENT,
		Port:      listenPort,
	}
	data, _ := json.Marshal(msg)
	broadConn.Write(data)
}
