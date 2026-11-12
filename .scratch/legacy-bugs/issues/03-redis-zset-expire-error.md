# zset add 子命令 Expire 错误被忽略

Status: resolved

## Problem

`cmd/client/redis.go` zSetAdd Run()（约 L405-407）：

```go
if expire > 0 {
    client.Expire(ctx, key, expire)
}
```

返回的 error 被忽略，过期设置失败时用户无感知。另外该 flag 描述写 "expire default 60s"，实际默认值是 0，文案与行为不符。

## Fix

检查 `client.Expire(...)` 的 error 并打印；修正 flag 描述为 "expire default 0s(不过期)"。

## Comments

- 2026-11-12: 评估时发现，未处理，保持 open。
- 2026-11-12: 已修复并在真实环境（192.168.1.10:6395, Redis Cluster）验证：`zset add --expire 120s` 两个 member 写入成功、`zset get` 按 score 倒序返回、`ttl` 确认 2m0s 生效。回归 build/vet/test 通过。
