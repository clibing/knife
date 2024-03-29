package system

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	str, ver, o, arch string
)

// archCmd represents the arch command
var archCmd = &cobra.Command{
	Use:   "arch",
	Short: "get current system kernal, machine.",
	Long: `get current os env:

默认替换os、arch关键字, 注意输入的-s的值需要使用单引号, 不能使用双信号。

变量支持 ${NAME} $NAME ${name} $name 四种定义方式，当使用$+变量名字时候，会出现替换错误， 例如 当OS=linux, 其中 $OSA和$OS, 会被替换 linuxA 和 linux
knife arch -s 'https://go.dev/dl/go${version}.${os}-${arch}.tar.gz' 
knife arch -s 'https://go.dev/dl/go${VERSION}.${OS}-$arch.tar.gz' 
knife arch -s 'https://go.dev/dl/go$version.$os-${ARCH}.tar.gz' 
knife arch -s 'https://go.dev/dl/go$VERSION.$OS-$ARCH.tar.gz' 

knife arch -s 'https://go.dev/dl/go${version}.${os}-${arch}.tar.gz' -o linux -a 386 -v 1.20.4

使用环境变量
set VV=99
knife arch -s 'https://go.dev/dl'"$VV"'/go${version}.${os}-${arch}.tar.gz' -o linux -a 386 -v 1.20.4

$GOOS      $GOARCH
darwin     386
darwin     amd64
darwin     arm
darwin     arm64
dragonfly  amd64
freebsd    386
freebsd    amd64
freebsd    arm
linux      386
linux      amd64
linux      arm
linux      arm64
linux      ppc64
linux      ppc64le
netbsd     386
netbsd     amd64
netbsd     arm
openbsd    386
openbsd    amd64
openbsd    arm
plan9      386
plan9      amd64
solaris    amd64
windows    386
windows    amd64

current os arch, kernal.`,
	Run: func(_ *cobra.Command, _ []string) {
		if len(str) == 0 {
			fmt.Printf("%s %s", runtime.GOOS, runtime.GOARCH)
		}

		result := str
		var Os, Architecture string
		if len(ver) > 0 {
			result = strings.ReplaceAll(result, "${version}", ver)
			result = strings.ReplaceAll(result, "$version", ver)
			result = strings.ReplaceAll(result, "${VERSION}", ver)
			result = strings.ReplaceAll(result, "$VERSION", ver)
		}

		if len(o) > 0 {
			Os = o
		} else {
			Os = runtime.GOOS
		}
		result = strings.ReplaceAll(result, "$OS", Os)
		result = strings.ReplaceAll(result, "${os}", Os)
		result = strings.ReplaceAll(result, "${OS}", Os)
		result = strings.ReplaceAll(result, "$os", Os)

		if len(arch) > 0 {
			Architecture = arch
		} else {
			Architecture = runtime.GOARCH
		}
		result = strings.ReplaceAll(result, "${arch}", Architecture)
		result = strings.ReplaceAll(result, "$arch", Architecture)
		result = strings.ReplaceAll(result, "${ARCH}", Architecture)
		result = strings.ReplaceAll(result, "$ARCH", Architecture)

		fmt.Println(result)
	},
}

func init() {
	archCmd.PersistentFlags().StringVarP(&str, "str", "s", "", "需要替换的URL")
	archCmd.PersistentFlags().StringVarP(&ver, "version", "v", "", "需要替换的版本号")
	archCmd.PersistentFlags().StringVarP(&o, "os", "o", "", "当前系统运行的操作系统")
	archCmd.PersistentFlags().StringVarP(&arch, "arch", "a", "", "当前系统运行的架构")
}

func NewArchCmd() *cobra.Command {
	return archCmd
}
