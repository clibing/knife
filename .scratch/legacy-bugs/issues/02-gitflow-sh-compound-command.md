# gitflow Install 的 sh 复合命令参数拆分错误

Status: resolved

## Problem

`cmd/system/pkg/gitflow_control.go` Install()：

```go
ExecuteCommand(value.Name, "sh", []string{"cd", value.Target, "&&", "make", "single"}, false)
```

`ExecuteCommand` 内部是 `exec.Command(bin, args...)`，实际执行 `sh cd <target> && make single`——sh 把 "cd" 当作要执行的脚本文件名（不存在，报错退出），`&&`/`make` 由外层 shell 解释。复合命令必须以 `sh -c "单字符串"` 形式传入。

## Fix

```go
ExecuteCommand(value.Name, "sh", []string{"-c", fmt.Sprintf("cd %s && make single", value.Target)}, false)
```

全仓 grep 确认无其他同类写法（ohmyzsh.go 的 `sh /tmp/ohmyzsh.sh` 是合法脚本文件调用）。

## Comments

- 2026-11-12: 已修复，`go build ./... && go vet ./... && go test ./...` 通过。
