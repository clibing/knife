package pkg

import (
	"log"
)

type Ohmyzsh struct{}

/**
 * 安装应用
 */
func (v *Ohmyzsh) Install(value *Package) bool {
	log.Printf("[%s]下载安装ohmyzsh安装脚本", value.Name)
	err := ExecuteCommand(value.Name, "curl", []string{"-fsSL", "https://install.ohmyz.sh/", "-o", "/tmp/ohmyzsh.sh"}, false)
	if err != nil {
		log.Printf("[%s]下载安装脚本异常%s", value.Name, err)
		return false
	}

	log.Printf("[%s]授权ohmyzsh可执行脚本", value.Name)
	// 授权
	ExecuteCommand(value.Name, "chmod", []string{"+x", "/tmp/ohmyzsh.sh"}, false)

	log.Printf("[%s]安装ohmyzsh", value.Name)
	err = ExecuteCommand(value.Name, "sh", []string{"/tmp/ohmyzsh.sh"}, false)
	if err != nil {
		log.Printf("[%s]安装失败%s", value.Name, err)
	}
	log.Printf("[%s]安装成功", value.Name)
	return true
}

/**
 * 升级应用
 */
func (v *Ohmyzsh) Upgrade(value *Package) bool {
	return true
}

/**
 * 之后事件
 */
func (v *Ohmyzsh) After(value *Package) {
}

func (v *Ohmyzsh) Before(value *Package, overwrite bool) bool {
	homeDir, _ := GetHomeDir(value.Name)
	has, _ := ExistPath(value.Name, homeDir, ".oh-my-zsh")
	return !has
}

func (v *Ohmyzsh) GetPackage() *Package {
	return &Package{
		Name:        "ohmyzsh",
		Bin:         "ohmyzsh",
		Version:     "latest",
		Env:         []*Env{},
		Shell:       ``,
		Compress:    "git",
		Target:      "",
		Description: "ohmyzsh源码安装",
		Source: []string{
			"-c",
			"$(curl -fsSL https://install.ohmyz.sh/)",
		},
	}
}
