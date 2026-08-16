# 修复前故障复现（Docker）
## 项目与标准命令
FieldOps Roster 是 Go 1.26 的本地工单排班 CLI。标准构建命令为 `go build ./...`，标准测试命令为 `go test ./...`。

## 环境构建与编译
在当前平台执行 `go build ./...` 可以完成编译，说明故障不是编译期问题。

## 故障触发步骤
1. 进入仓库根目录。
2. 执行 `go test ./...`。

## 实际错误输出
```text
--- FAIL: TestPlanSkipsInactiveTechniciansEvenWhenTheyMatchTheSkill (0.00s)
    inactive_test.go:25: inactive technician received work: [{OrderID:wo-fiber TechID:t-inactive ScheduledFor:2026-08-16 12:00:00 +0000 UTC Minutes:90}]
FAIL
```

## 期望行为
非活跃技师即使区域和技能匹配，也不能接收新工单；没有活跃匹配技师时，工单应保持未分配。
