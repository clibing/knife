package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/clibing/knife/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/net/webdav"
)

// webdavCmd represents the static command
var webdavCmd = &cobra.Command{
	Use:   "webdav",
	Short: "webdav服务。",
	Long: `启用本地webdav资源服务:

knife server webdav -p ./ 
`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			fmt.Println("webdav静态文件路径为空")
			return
		}
		wh := &webdav.Handler{
			FileSystem: webdav.Dir(path),
			LockSystem: webdav.NewMemLS(),
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			if username != "" && password != "" {
				authName, authPwd, ok := r.BasicAuth()
				if !ok {
					w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if username != authName || password != authPwd {
					http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
					return
				}
			}
			if r.Method == "GET" && handleDirList(wh.FileSystem, w, r) {
				return
			}
			readonly, _ := cmd.Flags().GetBool("readonly")
			if readonly {
				switch r.Method {
				case "PUT", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE":
					http.Error(w, "WebDAV: Read Only!!!", http.StatusForbidden)
					return
				}
			}
			wh.ServeHTTP(w, r)
		})

		targetPort, err := cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Println("启动端口异常: ", targetPort)
			return
		}

		// 获取端口
		value := fmt.Sprintf(":%d", targetPort)
		// 建立监听
		listener, err := net.Listen("tcp", value)
		if err != nil {
			fmt.Println("建立监听异常, ", err.Error())
			return
		}

		// 获取 本地ip地址
		ip, _ := utils.GetLocalIp()
		latestPort := listener.Addr().(*net.TCPAddr).Port

		fmt.Printf("服务启动中: http://%s:%d \n", ip, latestPort)
		fmt.Printf("服务启动中: http://0.0.0.0:%d \n", latestPort)
		err = http.Serve(listener, nil)
		if err != nil {
			fmt.Println("服务启动失败，请检查, ", err.Error())
			return
		}
	},
}

func handleDirList(fs webdav.FileSystem, w http.ResponseWriter, req *http.Request) bool {
	ctx := context.Background()
	f, err := fs.OpenFile(ctx, req.URL.Path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()
	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {
		return false
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(w, "Error reading directory", http.StatusInternalServerError)
		return false
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", name, name)
	}
	fmt.Fprintf(w, "</pre>\n")
	return true
}

func init() {
	webdavCmd.Flags().StringP("path", "p", "./", "静态文件路径")
	webdavCmd.Flags().IntP("port", "", 0, "端口, 默认会随机")
	webdavCmd.Flags().StringP("username", "u", "guest", "webdav basic auth username")
	webdavCmd.Flags().StringP("password", "P", "guest", "webdav basic auth password")
	webdavCmd.Flags().BoolP("readonly", "r", false, "webdav启用只读模式")
}

func NewWebdavServer() *cobra.Command {
	return webdavCmd
}
