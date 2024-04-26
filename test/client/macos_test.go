package client

import (
	"testing"

	"github.com/clibing/knife/cmd/system/pkg"
)

func TestMacos(t *testing.T) {
	ohmyzsh := pkg.NewOhmyzshPlugin()
	p := ohmyzsh.GetPackage()
	ohmyzsh.Install(p)
}

func TestMacos2(t *testing.T) {
	b := &pkg.Brew{}
	p := b.GetPackage()
	b.Install(p)
}