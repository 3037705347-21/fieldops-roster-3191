# fieldops-roster__002 Docker 交付说明

## 项目概览
- FieldOps Roster is a local command-line tool for small field-service teams that need to schedule work orders without an external dispatch platform. Dispatch coordinators can valida
- Go module: `fieldops-roster`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/rosterctl
```

## Docker 构建

```bash
./build_benzhi_docker.sh fieldops-roster__002-benzhi linux/amd64
docker run --rm -it fieldops-roster__002-benzhi bash
```

## 环境

- 基础镜像: `golang:1.26`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
