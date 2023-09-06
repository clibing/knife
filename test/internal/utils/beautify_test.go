package utils

import (
	"fmt"
	"testing"

	"github.com/clibing/knife/internal/utils"
)

func TestBeautify(t *testing.T) {
	size := int64(2149)
	value := utils.BeautifyValue(size)
	fmt.Println(value)

	u, v := utils.BeautifyUnit(size)
	fmt.Println(u, v)
}
