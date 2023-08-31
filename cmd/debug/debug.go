package debug

import (
	"fmt"

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
	return &Debug{
		Command:         c,
		CombinationKey:  true,
		Mode:            RELEASE,
		AutoAppendEntry: true,
	}
}

func (d *Debug) EnableDebug() (debug bool) {
	debug, e := d.Command.Flags().GetBool(fmt.Sprintf("%s-debug", d.Command.Use))
	if e != nil {
		debug = false
	}
	return
}

func (d *Debug) Show(developFormat, releaseFormat string, parameters ...interface{}) {
	if d.EnableDebug() {
		fmt.Printf(d.AppendEntry(developFormat), parameters...)
	} else {
		fmt.Printf(d.AppendEntry(releaseFormat), parameters...)
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
