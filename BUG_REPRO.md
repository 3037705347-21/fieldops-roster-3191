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
--- FAIL: TestPlanPrioritizesUrgentDueOrdersAndBalancesCapacity (0.00s)
    planner_test.go:30: first assignment = wo-low, want urgent order first
--- FAIL: TestPlanLeavesOrdersUnassignedWhenSkillOrCapacityIsMissing (0.00s)
    planner_test.go:53: assignments = 2, want 1
--- FAIL: TestAssignOpenOrdersUpdatesWorkspaceState (0.00s)
    dispatch_test.go:30: order not scheduled correctly: {ID:wo-1 Customer:Clinic A Region:north RequiredSkill:fiber DurationMinutes:60 Priority:urgent Status:scheduled DueAt:2026-08-16 10:00:00 +0000 UTC AssignedTechID: CompletedAt:<nil>}
FAIL
```

## 期望行为
紧急工单应优先排班，技师容量不能超出当天上限，排班保存后工单应能看到实际指派的技师。
