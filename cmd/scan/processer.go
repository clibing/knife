package scan

import (
	"fmt"
	"io/fs"
	"time"
)

type Java struct{}
type Rust struct{}
type Golang struct{}
type C struct{}
type Python struct{}

func (j *Java) Processer(path string, info fs.FileInfo) (meta *Meta, ok bool, err error) {
	time.Sleep(1 * time.Second)
	return
}

func (v *Java) Log(format string, args ...interface{}) {
	fmt.Printf("[Java] "+format+"\n", args...)
}

func (r *Rust) Processer(path string, info fs.FileInfo) (meta *Meta, ok bool, err error) {
	time.Sleep(1 * time.Second)
	return
}

func (v *Rust) Log(format string, args ...interface{}) {
	fmt.Printf("[Rust] "+format+"\n", args...)
}

func (g *Golang) Processer(path string, info fs.FileInfo) (meta *Meta, ok bool, err error) {
	time.Sleep(1 * time.Second)
	return
}

func (v *Golang) Log(format string, args ...interface{}) {
	fmt.Printf("[Golang] "+format+"\n", args...)
}

func (c *C) Processer(path string, info fs.FileInfo) (meta *Meta, ok bool, err error) {
	time.Sleep(1 * time.Second)
	return
}

func (v *C) Log(format string, args ...interface{}) {
	fmt.Printf("[C] "+format+"\n", args...)
}

func (p *Python) Processer(path string, info fs.FileInfo) (meta *Meta, ok bool, err error) {
	time.Sleep(1 * time.Second)
	return
}

func (v *Python) Log(format string, args ...interface{}) {
	fmt.Printf("[Python] "+format+"\n", args...)
}
