package pkg

import (
	"fmt"
	"log"
	"os/user"
)

type Go struct{}

/**
 * 安装应用
 */
func (v *Go) Install(value *Package) bool {
	has, e := CheckCommand(value.Name, value.Bin)
	if !has {
		log.Printf("[%s]暂未安装, has: %t, err: %s\n", value.Name, has, e)

	}
	result := FilterAppendEnv(value.Name, value.Env)
	status, e := SettingEnv(value.Name, result, value.Target)
	if e != nil {
		log.Printf("[%s]%s", value.Name, e.Error())
		return false
	}
	return status
}

/**
 * 升级应用
 */
func (v *Go) Upgrade(value *Package) bool {
	return true
}

/**
 * 前置事件
 */
func (v *Go) Before(value *Package, overwrite bool) bool {
	profile, e := GetCurrentProfile(value.Name)
	if e != nil {
		log.Printf("[%s]安装Go环境变量, 错误信息:%s", value.Name, e.Error())
		return false
	}

	log.Printf("[%s]向文件:%s, 追加内容", value.Name, profile)
	value.Target = profile
	return true
}

/**
 * 后置事件
 */
func (v *Go) After(value *Package) {
	log.Printf("[%s]Go环境变量完成", value.Name)
}

func (v *Go) GetPackage() *Package {
	u, e := user.Current()
	if e != nil {
		log.Println("获取当前用户名异常", e)
	}

	return &Package{
		Name:    v.Key(),
		Bin:     "go",
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


func (v *Go) Key() string {
	return "Go"
}
