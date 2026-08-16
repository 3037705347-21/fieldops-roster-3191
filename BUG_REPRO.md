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
FAIL
```

## 期望行为
取消已排班工单后，旧排班记录应同步移除，后续重新排班时技师容量应按最新工作区重新计算。
