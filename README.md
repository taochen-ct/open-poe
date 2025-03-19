# Gin 项目模板

本项目是一个基于 Gin 框架的 Golang Web API 服务，集成了 `wire` 依赖注入、`cobra` 命令行管理、`zap` 日志库，并结合 `lumberjack` 进行日志切割。

## 主要特性
- **Gin**：高性能 HTTP Web 框架。
- **Wire**：依赖注入工具，简化对象创建和管理。
- **Cobra**：命令行管理工具，支持 CLI 结构化管理。
- **Zap**：高性能日志库。
- **Lumberjack**：日志切割与归档管理。

## 依赖安装

确保你已安装 Go，并在项目根目录执行：

```sh
# 初始化 go module
go mod tidy
```

## 代码结构

```
.
├── cmd                 # CLI 命令管理
│   ├── api
│       ├── app.go      # API 相关命令
│       ├── main.go     # 入口文件
│       ├── wire.go     # Wire 依赖注入
├── config              # 配置声明
├── internal            # 内部业务逻辑
├── pkg                 # 通用组件
├── routes              # 注册路由
├── storage             # 存储目录
├── go.mod              # Go 依赖管理
├── conf                # yaml配置文件
```

## 启动项目

```sh
# 生成 wire 依赖
wire

# 运行 API 服务
make run

# 或手动运行
go run ./cmd/api
```

## 使用 Cobra 运行 CLI

```sh
# 查看所有可用命令
./main --help

# 运行指定命令
./main example
```

## 配置日志
本项目使用 `zap` 结合 `lumberjack` 进行日志管理，配置示例：

```go
logger := zap.NewProduction()
logger = logger.WithOptions(zap.Hooks(lumberjack.Logger{
    Filename:   "logs/app.log",
    MaxSize:    10, // MB
    MaxBackups: 3,
    MaxAge:     28, // days
}))
```


## 许可证
MIT License

