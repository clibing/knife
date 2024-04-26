package pkg

import (
	"fmt"
	"log"
)

type Ohmyzsh struct{}

/**
 * 安装应用
 */
func (v *Ohmyzsh) Install(value *Package) bool {
	log.Printf("[%s]下载安装ohmyzsh安装脚本\n", value.Name)
	err := ExecuteCommand(value.Name, "curl", []string{"-fsSL", "https://install.ohmyz.sh/", "-o", "/tmp/ohmyzsh.sh"}, false)
	if err != nil {
		log.Printf("[%s]下载安装脚本异常%s\n", value.Name, err)
		return false
	}

	log.Printf("[%s]授权ohmyzsh可执行脚本\n", value.Name)
	// 授权
	ExecuteCommand(value.Name, "chmod", []string{"+x", "/tmp/ohmyzsh.sh"}, false)

	log.Printf("[%s]安装ohmyzsh\n", value.Name)
	err = ExecuteCommand(value.Name, "sh", []string{"/tmp/ohmyzsh.sh"}, false)
	if err != nil {
		log.Printf("[%s]安装失败%s\n", value.Name, err)
	}
	log.Printf("[%s]安装成功\n", value.Name)
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
	if overwrite {
		log.Printf("[%s]强制安装\n", value.Name)
		ExecuteCommand(value.Name, "rm", []string{"-rf", "/tmp/ohmyzsh.sh"}, false)
		ExecuteCommand(value.Name, "rm", []string{"-rf", fmt.Sprintf("%s/.oh-my-zsh", homeDir)}, false)
		return true
	}
	if has {
		log.Printf("[%s]已安装\n", value.Name)
	}
	return !has
}

func (v *Ohmyzsh) GetPackage() *Package {
	return &Package{
		Name:        v.Key(),
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

func (v *Ohmyzsh) Key() string {
	return "ohmyzsh"
}
