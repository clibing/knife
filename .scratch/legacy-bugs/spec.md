# legacy-bugs

修复 2026-11-12 架构评估报告中遗留的两个未验证 bug，并把仍未处理的观察项登记为 open issue。

## Goals

- 修复 `cmd/client/redis.go` SET 分支 expire 解析条件反转（报告 redis.go:114-117）
- 修复 `cmd/system/pkg/gitflow_control.go` Install 中 sh 复合命令参数拆分错误（报告 gitflow_control.go:15）
- 登记评估过程中发现但未处理的次要问题

## Non-goals

- 不重构 redis.go / gitflow_control.go 的整体结构
- 不为需要真实 Redis 连接的路径补集成测试（已有 `//go:build integration` 约定）

## Status

已完成：01-04 全部修复，并在真实环境（192.168.1.10:6395 Cluster + 6379 单实例）双环境实测通过。`-c` 模式定位为单实例/主从场景（作者确认），cluster 用户走子命令。后续可选项：失败路径非 0 exit code（见 issue 04 未做部分）。

## Issues

- issues/01-redis-expire-condition.md
- issues/02-gitflow-sh-compound-command.md
- issues/03-redis-zset-expire-error.md
- issues/04-redis-top-level-cmd-cluster-silent.md
