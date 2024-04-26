package pkg

import (
	"fmt"
	"os/user"
)

type Golang struct{ Log }

/**
 * 安装应用
 */
func (v *Golang) Install(value *Package) bool {
	has, e := CheckCommand(value.Name, value.Bin)
	if !has {
		v.Log.Println("暂未安装, has: %t, err: %s", has, e)
	}
	if !has {
		v.Log.Println("安装%s", value.Bin)
		InstallByBrew(value.Name, value.Bin)
	}

	result := FilterAppendEnv(value.Name, value.Env)
	status, e := SettingEnv(value.Name, result, value.Target)
	if e != nil {
		v.Log.Println("%s", e.Error())
		return false
	}
	return status
}

/**
 * 升级应用
 */
func (v *Golang) Upgrade(value *Package) bool {
	return true
}

/**
 * 前置事件
 */
func (v *Golang) Before(value *Package, overwrite bool) bool {
	profile, e := GetCurrentProfile(value.Name)
	if e != nil {
		v.Log.Println("检查用户配置异常，错误信息:%s", e.Error())
		return false
	}

	v.Log.Println("开始安装:%s, 追加内容", profile)
	value.Target = profile
	return true
}

/**
 * 后置事件
 */
func (v *Golang) After(value *Package) {
	v.Log.Println("安装完成")
}

func (v *Golang) GetPackage() *Package {
	u, e := user.Current()
	if e != nil {
		v.Log.Println("获取当前用户名异常:%s", e)
	}

	return &Package{
		Name:    v.Key(),
		Bin:     "golang",
		Version: "latest",
		Env: []*Env{
			{
				Key:   "GO111MODULE",
				Value: "auto",
			},
			{
				Key:   "GOPROXY",
				Value: "https://goproxy.cn,https://mirrors.aliyun.com/goproxy,direct",
			},
			{
				Key:   "GOPATH",
				Value: fmt.Sprintf("%s/go", u.HomeDir),
			},
			{
				Key:   "GOBIN",
				Value: "$GOPATH/bin",
			},
		},
		Target:      "",
		Description: "Go环境变量",
		Source:      []string{},
	}
}

func (v *Golang) Key() string {
	return "golang"
}

func NewGolang() *Golang {
	g := &Golang{}
	l := Log{Key: g.Key()}
	g.Log = l
	return g
}
