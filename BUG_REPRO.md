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
--- FAIL: TestCancelOrderRemovesAssignment (0.00s)
    dispatch_test.go:57: assignments = [{OrderID:wo-1 TechID:t-1 ScheduledFor:2026-08-16 10:00:00 +0000 UTC Minutes:60}], want empty
--- FAIL: TestCompleteOrderRejectsOpenOrders (0.00s)
    dispatch_test.go:72: error = <nil>, want ErrInvalidTransition
--- FAIL: TestValidateWorkspaceRejectsAssignmentForUnknownOrder (0.00s)
    jsonstore_test.go:27: error = unknown work order: wo-missing, want ErrUnknownOrder
FAIL
```

## 期望行为
open 工单不能跳过处理过程直接完成，取消已排班工单后不应继续保留容量占用，校验错误应保留可识别的业务错误类型。
