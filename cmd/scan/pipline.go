package scan

import (
	"fmt"
	"time"
)

type Sqlite3 struct {
}

type Console struct {
}

func (s *Sqlite3) Pipline(meta *Meta) (err error) {
	s.Log("pipline wiht sqlite3")
	time.Sleep(1*time.Second)
	return
}

func (ds *Sqlite3) Log(format string, args ...interface{}) {
	fmt.Printf("[Sqlite3] "+format+"\n", args...)
}

func (c *Console) Pipline(meta *Meta) (err error) {
	c.Log("pipline wiht console")
	time.Sleep(1*time.Second)
	return
}

func (ds *Console) Log(format string, args ...interface{}) {
	fmt.Printf("[Console] "+format+"\n", args...)
}
