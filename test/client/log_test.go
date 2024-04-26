package client

import (
	"testing"

	"github.com/clibing/knife/cmd/system/pkg"
)

func TestLogFormat(t *testing.T) {
	p := []interface{}{"value"}
	pkg.LogFormatted("内容: %s", "test", p)
}
