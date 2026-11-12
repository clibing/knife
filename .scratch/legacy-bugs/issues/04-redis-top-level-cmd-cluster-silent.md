# 顶层 -c 模式（GET/SET/DEL）的静默行为

Status: resolved

## 定位（2026-11-12 与作者确认）

`knife client redis -c GET/SET/DEL`（redisCmd，单机 `redis.NewClient` + `master()` 判断）面向**单实例/主从 sentinel** 场景（如 192.168.1.10:6379，实测 role:master、connected_slaves:0）。Redis Cluster（6391-6396）不是它的目标场景，cluster 用户应走子命令（get/set/ttl/del/zset，走 ClusterClient，MOVED 自动跟随）。

## 单实例实测结果（192.168.1.10:6379，2026-11-12，修复前）

正常路径全部通过：

- `SET -k knife-agent-k1 -v v1 -e 60s` → `OK. set knife-agent-k1 v1`
- `GET -k knife-agent-k1` → `  [knife-agent-k1]-->[v1] expire: 60s`
- `DEL -k knife-agent-k1` → 实际删除成功（exists=0 复核）
- 多 host 遍历 + role 识别：`0: 192.168.1.10:6379 role(master)` / `1: 192.168.1.10:6395 role(slave)`
- 非法 expire `-e abc` → `expire format error: abc` 并 abort（issue 01 修复路径）
- `-e ""` 被 cobra 解析为空 slice → 走 "key, value, expire丢失" 分支（L114 空串分支为防御路径，行为无回归）

## 已修复（2026-11-12）

`execute()` 三处静默行为改为明确输出：

1. **GET**：区分 `redis.Nil`（打印 `key (...) does not exist`，继续下一 key）与其他错误（打印错误 + abort，中断本节点）。
2. **DEL**：打印每个 key 的删除数量（`del <key>: <n>`）；删除失败打印错误 + abort；非 master 打印 `当前节点非master, 跳过DEL`。
3. **SET**：非 master 打印 `当前节点非master, 跳过SET`。

注：`-k k1 k2` 只消费 k1（cobra StringSliceVar 语义），多 key 需 `-k k1 -k k2`，本次未改动该 flag 语义。

## 修复后实测（2026-11-12）

- 6379 单实例：SET/GET 正常输出不变；`GET knife-agent-missing` → `key (knife-agent-missing) does not exist`；`DEL` → `del knife-agent-k1: 1`
- 6395 slave：`SET` → `当前节点非master, 跳过SET`；`DEL` → `当前节点非master, 跳过DEL`
- `go vet ./... && go test ./...` 全绿

## 未做（记录为后续可选项）

- 关键失败路径返回非 0 exit code：需将 `Run` 改为 `RunE` 并向上传播，影响整个 CLI 行为契约，暂不动。

## Comments

- 2026-11-12: cluster 环境验证 issue 01/03 时发现；与作者确认 6379 为单实例后修正定位，单实例正向路径实测通过。
- 2026-11-12: 三处静默行为已修复并双环境实测通过，转 resolved。
