package pkg

import "fmt"

type Rust struct {
	Log
}

func (v *Rust) Install(value *Package) bool {
	e := ExecuteCommand(value.Name, "curl", []string{"-fsSL", value.Source[0], "-o", value.Shell}, false)
	if e != nil {
		v.Log.Println("下载Rust安装脚本: %s", e)
		return false
	}

	e = ExecuteCommand(value.Name, "chmod", []string{"+x", value.Shell}, false)
	if e != nil {
		v.Log.Println("授权安装脚本可执行: %s", e)
		return false
	}

	e = ExecuteCommand(value.Name, "sh", []string{value.Shell}, false)
	if e != nil {
		v.Log.Println("安装错误: %s", e)
		return false
	}

	return true
}

func (v *Rust) Upgrade(value *Package) bool {
	homeDir, e := GetCurrentProfile(value.Name)
	if e != nil {
		return false
	}
	bin := fmt.Sprintf("%s/.cargo/bin/rustup", homeDir)
	ExecuteCommand(value.Name, bin, []string{"update"}, false)

	return true
}

func (v *Rust) Before(value *Package, overwrite bool) bool {
	has, _ := CheckCommand(value.Name, "rustup")
	// if e != nil {
	// 	v.Log.Println("Check rustup failed: %s", e)
	// 	return
	// }
	e := ExecuteCommand(value.Name, "rm", []string{"-rf", value.Shell}, false)
	if e != nil {
		v.Log.Println("删除文件: %s发生错误: %s", value.Shell, e)
		return false
	}
	return !has
}

func (v *Rust) After(value *Package) {
	v.Log.Println("安装完成")

}

func (v *Rust) GetPackage() *Package {
	return &Package{
		Name:        "Rust",
		Bin:         "rustup",
		Shell:       "/tmp/rust-install.sh",
		Version:     "latest",
		Env:         []*Env{},
		Description: "Rust安装",
		Source: []string{
			"https://sh.rustup.rs",
		},
	}
}

func (v *Rust) Key() string {
	return "Rust"
}

func NewRust() *Rust {
	tmp := &Rust{}
	l := Log{Key: tmp.Key()}
	tmp.Log = l
	return tmp
}
