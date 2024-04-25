package client

import (
	"testing"

	"github.com/clibing/knife/cmd/system/pkg"
)

func TestMacos(t *testing.T) {
	ohmyzsh := pkg.NewOhmyzsh()
	p := ohmyzsh.GetPackage()
	ohmyzsh.Install(p)
}
