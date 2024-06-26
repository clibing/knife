package pkg

import (
	"fmt"
	"log"
	"os"
)

type GitflowControl struct {
	Log
}

func (v *GitflowControl) Install(value *Package) bool {
	ExecuteCommand(value.Name, "knife", []string{"client", "goc", "-s", value.Source[0]}, false)
	ExecuteCommand(value.Name, "sh", []string{"cd", value.Target, "&&", "make", "single"}, false)
	return true
}

func (v *GitflowControl) Upgrade(value *Package) bool {
	return true
}

func (v *GitflowControl) Before(value *Package, overwrite bool) bool {
	has, _ := CheckCommand(value.Name, "knife")
	if !has {
		log.Printf("[%s]请先安装本工具knife, https://github.com/clibing/knife/releases.\n", value.Name)
		return has
	}
	return true
}

func (v *GitflowControl) After(value *Package) {
	log.Printf("[%s]安装完成\n", value.Name)
}

func (v *GitflowControl) Key() string {
	return "gitflow-control"
}

func (v *GitflowControl) GetPackage() *Package {
	gopath := os.Getenv("GOPATH")
	return &Package{
		Name:        v.Key(),
		Bin:         "gitflow-control",
		Version:     "latest",
		Env:         []*Env{},
		Shell:       "",
		Compress:    "binary",
		Target:      fmt.Sprintf("%s/src/github.com/clibing/gitflow-control", gopath),
		Description: "git扩展功能, 例如: git ci, git add, git issue.",
		Source: []string{
			"https://github.com/clibing/gitflow-control.git",
		},
	}
}

func NewGitflowControl() *GitflowControl {
	g := &GitflowControl{}
	g.Log = Log{Key: g.Key()}
	return g
}
