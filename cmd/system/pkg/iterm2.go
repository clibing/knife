package pkg

import "log"

type ITerm2 struct {
	Log
}

func (v *ITerm2) Install(value *Package) bool {
	tmp := "/tmp/iterm2.zip"

	url := value.Source[0]
	err := ExecuteCommand(value.Name, "curl", []string{"-fsSL", "-o", tmp, url}, false)
	if err != nil {
		log.Printf("[%s]下载文件异常\n", value.Name)
		return false
	}

	err = ExecuteCommand(value.Name, "unzip", []string{tmp, "-d", value.Target}, false)
	if err != nil {
		log.Printf("[%s]解药到%s失败\n", value.Name, err.Error())
		return false
	}
	return true
}

func (v *ITerm2) Upgrade(value *Package) bool {
	return true
}

func (v *ITerm2) Before(value *Package, overwrite bool) bool {
	has, e := ExistApplications(value.Name, value.Bin)
	if e != nil {
		log.Printf("[%s]检查应用异常: %s\n", value.Name, e.Error())
	}
	if overwrite {
		log.Printf("[%s]强制安装\n", value.Name)
		return true
	}
	// 不存在需要执行安装
	return !has
}

func (v *ITerm2) After(value *Package) {

}

func (v *ITerm2) GetPackage() *Package {
	return &Package{
		Name:        v.Key(),
		Bin:         "iTerm.app",
		Version:     "latest",
		Compress:    "zip",
		Target:      "/Applications",
		Description: "iTerm2是一个开源的 macOS 终端模拟器。",
		Source:      []string{"https://iterm2.com/downloads/stable/latest"},
	}
}

func (v *ITerm2) Key() string {
	return "iTerm2"
}

func NewITerm2() *ITerm2 {
	i := &ITerm2{}
	l := Log{Key: i.Key()}
	i.Log = l
	return i

}
