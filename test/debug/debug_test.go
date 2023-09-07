package debug

import (
	"fmt"
	"testing"

	"github.com/clibing/knife/cmd/debug"
	"github.com/spf13/cobra"
)

func TestFillSpace(t *testing.T) {
	c := &cobra.Command{
		Use: "ip",
	}
	debug := debug.NewDebug(c)

	var value string
	middle := " : "
	value = debug.FillSpace(10, middle, "admin")
	fmt.Println(value)
	value = debug.FillSpace(10, middle, "clibing")
	fmt.Println(value)
	value = debug.FillSpace(10, middle, "wmsjhappy")
	fmt.Println(value)
	value = debug.FillSpace(10, middle, "name")
	fmt.Println(value)
	value = debug.FillSpace(10, middle, "abc")
	fmt.Println(value)
}
