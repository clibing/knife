package system

import (
	"fmt"
	"testing"

	"github.com/clibing/knife/cmd/system/pkg"
)

func TestMacosCheckCmd(t *testing.T) {
	has, e := pkg.CheckCommand("test", "DEMO")
	fmt.Println("DEMO", has, e)
	has, e = pkg.CheckCommand("test", "java")
	fmt.Println("java", has, e)
	has, e = pkg.CheckCommand("test","go")
	fmt.Println("go", has, e)
}
