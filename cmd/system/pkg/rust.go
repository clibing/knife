package pkg

type Rust struct {
	Log
}

func (v *Rust) Install(value *Package) bool {
	ExecuteCommand(value.Name, "curl", []string{"-o", value.Shell, "--proto", "'=https'", "--tlsv1.2", "-sSf", value.Source[0]}, false)
	ExecuteCommand(value.Name, "chmod", []string{"+x", value.Shell}, false)
	ExecuteCommand(value.Name, "sh", []string{value.Shell}, false)
	return true
}

func (v *Rust) Upgrade(value *Package) bool {
	return true
}

func (v *Rust) Before(value *Package, overwrite bool) bool {
	has, _ := CheckCommand(value.Name, "rustup")
	ExecuteCommand(value.Name, "rm", []string{"-rf", "/tmp/rust-install.sh"}, false)
	return has
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
