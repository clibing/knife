package client

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// gocCmd represents the random command
var gocCmd = &cobra.Command{
	Use:   "goc",
	Short: "git clone golang project to $GOPATH",
	Long: `git clone golang project:

1. goc -s https://github.com/clibing/knife.git.
2. goc -s https://github.com/clibing/knife.git -b master.
`,

	Run: func(cmd *cobra.Command, arg []string) {
		source, _ := cmd.Flags().GetString("source")
		branch, _ := cmd.Flags().GetString("branch")
		output, _ := cmd.Flags().GetString("output")

		if len(source) > 0 {
			fmt.Println("source:", source)
		} else {
			fmt.Println("请输入Git source")
			return
		}

		var parameters []string
		parameters = append(parameters, "clone", source)

		if len(branch) > 0 {
			fmt.Println("branch:", source)
			parameters = append(parameters, "-b", branch)
		}

		var e error
		if len(output) == 0 {
			output, e = createDir(source)
			if e != nil {
				fmt.Println("生成项目目录错误", e)
				return
			}
		}

		fmt.Println("path  :", output)

		parameters = append(parameters, output)

		var result string
		// result, e = utils.ExecGit(parameters...)
		if e != nil {
			fmt.Println("faild :", result, e)
			return
		}

		if len(result) > 0 {
			fmt.Println("Result:", result)
		}

		time.Sleep(3 * time.Second)

		fmt.Printf("cmd   : git clone %s %s\n", source, output)
	},
}

func init() {
	gocCmd.Flags().StringP("source", "s", "", "git project url")
	gocCmd.Flags().StringP("branch", "b", "", "git branch name")
	gocCmd.Flags().StringP("output", "o", "", "将项目clone到指定目录")
}

func NewGocCmd() *cobra.Command {
	return gocCmd
}

func createDir(source string) (value string, e error) {
	gopath := os.Getenv("GOPATH")
	fmt.Println("GOPATH:", gopath)
	if len(source) == 0 {
		e = fmt.Errorf("source is empty")
		return
	}

	u, e := url.Parse(source)
	if e != nil {
		return
	}

	tmp, _ := strings.CutPrefix(u.Path, "/")

	tmp, _ = strings.CutSuffix(tmp, ".git")

	value = filepath.Join(gopath, "src", u.Host, tmp)
	return
}
