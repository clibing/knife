package debug

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

type Mode int

const (
	DEVELOP Mode = 1 // 开发模式
	RELEASE Mode = 2 // 发行模式
)

type Debug struct {
	Command         *cobra.Command // cobra command
	CombinationKey  bool           // debug flag key is combination key, default true. Command.Use+"-debug"
	Mode            Mode           // 运行模式
	AutoAppendEntry bool           // auto append entry

}

func NewDebug(c *cobra.Command) *Debug {
	d := &Debug{
		Command:         c,
		CombinationKey:  true,
		Mode:            RELEASE,
		AutoAppendEntry: true,
	}
	if d.EnableDebug() {
		fmt.Println("🟣 Enable debug")
		d.Mode = DEVELOP
	}
	return d
}

func (d *Debug) EnableDebug() (debug bool) {
	if d.Command == nil {
		return false
	}
	debug, e := d.Command.Flags().GetBool("debug")
	if e == nil {
		return
	}
	debug = false
	return
}

func (d *Debug) Show(developFormat, releaseFormat string, parameters ...interface{}) {
	if d.Mode == DEVELOP {
		if len(developFormat) > 0 {
			fmt.Printf(d.AppendEntry(developFormat), parameters...)
		}
	} else {
		if len(releaseFormat) > 0 {
			fmt.Printf(d.AppendEntry(releaseFormat), parameters...)
		}
	}
}

func (d *Debug) Debug(developFormat string, parameters ...interface{}) {
	d.Show(developFormat, "", parameters...)
}

func (d *Debug) Release(releaseFormat string, parameters ...interface{}) {
	d.Show("", releaseFormat, parameters...)
}

func (d *Debug) AppendEntry(source string) string {
	if d.AutoAppendEntry {
		return source + "\n"
	}
	return source
}

func (d *Debug) ShowSame(format string, parameters ...interface{}) {
	d.Show(format, format, parameters...)
}

/**
 * 填充空格
 */
func (d *Debug) FillSpace(width int, middle, value string) string {
	return fmt.Sprintf("%"+strconv.Itoa(width)+"s"+middle, value)
}
