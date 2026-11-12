# redis SET 分支 expire 解析条件反转

Status: resolved

## Problem

`cmd/client/redis.go` execute() 的 SET 分支：

```go
if len(expire) == 0 {
    expireDuration, _ = time.ParseDuration(expire) // 解析空串，错误被丢弃
} else {
    expireDuration, err = time.ParseDuration(expire)
    ...
}
```

空串分支去解析空串（必然报错且被 `_` 丢弃），非空分支才做错误处理——条件与意图相反。功能上碰巧成立（空串解析失败得到 0 时长 = 不过期），但结构反转、意图不清。

## Fix

条件改为 `if len(expire) > 0`：非空时解析并校验错误（abort），空串时直接保持 0 时长（不过期）。行为不变，结构归位。

## Comments

- 2026-11-12: 已修复，`go build ./... && go vet ./... && go test ./...` 通过。
