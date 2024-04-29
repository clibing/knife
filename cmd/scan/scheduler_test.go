package scan

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	bg := context.Background()
	ds := &DefaultScheduler{
		Done: make(chan bool, 1),
	}
	ds.BindProcesser(&Java{}, &Python{}, &C{}, &Golang{}, &Rust{})
	ds.BindPipline(&Sqlite3{}, &Console{})

	// ctx, _ := context.WithCancel(bg)
	// go
	ds.Start(bg, "./")

	// 等待
	var done bool
	for {
		select {
		case done = <-ds.Done:
			fmt.Println("调度器结束", done)
			goto End
		case <-time.After(1 * time.Second):
			fmt.Print(".")
		}
	}
End:
	ds.Log("End: %t", done)
}
