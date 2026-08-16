# 修复前故障复现（Docker）
## 项目与标准命令
FieldOps Roster 是 Go 1.26 的本地工单排班 CLI。标准构建命令为 `go build ./...`，标准测试命令为 `go test ./...`，排班入口为 `go run ./cmd/rosterctl plan -file testdata/workspace.json`。

## 环境构建与编译
在当前平台执行 `go build ./...` 可以完成编译，说明故障不是编译期问题。

## 故障触发步骤
1. 进入仓库根目录。
2. 执行 `go run ./cmd/rosterctl plan -file testdata/workspace.json`。

## 实际错误输出
```text
unknown technician: wo-1001
exit status 1
```

## 期望行为
plan 命令应把样例工作区中的 open 工单排给匹配技师并保存，不应在保存阶段报未知技师或未知工单。
