package pkg

import (
	"fmt"
	"log"
)

type Ohmyzsh struct {
	Log
}

/**
 * 安装应用
 */
func (v *Ohmyzsh) Install(value *Package) bool {
	v.Log.Println("下载安装ohmyzsh安装脚本")
	err := ExecuteCommand(value.Name, "curl", []string{"-fsSL", "https://install.ohmyz.sh/", "-o", "/tmp/ohmyzsh.sh"}, false)
	if err != nil {
		v.Log.Println("下载安装脚本异常%s", err)
		return false
	}

	v.Log.Println("授权ohmyzsh可执行脚本")
	// 授权
	ExecuteCommand(value.Name, "chmod", []string{"+x", "/tmp/ohmyzsh.sh"}, false)

	v.Log.Println("安装ohmyzsh")
	err = ExecuteCommand(value.Name, "sh", []string{"/tmp/ohmyzsh.sh"}, false)
	if err != nil {
		v.Log.Println("安装失败%s", err.Error())
	}
	v.Log.Println("安装成功")
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
		v.Log.Println("强制安装")
		ExecuteCommand(value.Name, "rm", []string{"-rf", "/tmp/ohmyzsh.sh"}, false)
		ExecuteCommand(value.Name, "rm", []string{"-rf", fmt.Sprintf("%s/.oh-my-zsh", homeDir)}, false)
		return true
	}
	if has {
		v.Log.Println("已安装")
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

func NewOhmyzsh() *Ohmyzsh {
	o := &Ohmyzsh{}
	l := Log{Key: o.Key()}
	o.Log = l
	return o
}
