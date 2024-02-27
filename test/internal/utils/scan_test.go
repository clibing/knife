package utils

import (
	"fmt"
	"testing"

	"github.com/clibing/knife/cmd/debug"
	"github.com/clibing/knife/internal/utils"
)

func TestScan(t *testing.T) {
	d := debug.NewDebug(nil)
	result := utils.Scan(d, "/Users/clibing/Downloads/英文")
	for _, value := range result {
		fmt.Println(value)
	}
}
