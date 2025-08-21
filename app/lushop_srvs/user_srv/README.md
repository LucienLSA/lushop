# user_srv 服务说明

## 项目简介

`user_srv` 是 lushop 微服务架构中的用户服务，负责用户相关的业务逻辑处理，包括用户数据的增删改查、认证、注册等功能。该服务基于 Go 语言开发，采用 gRPC 通信协议，并支持 Consul 服务注册与健康检查。

## 目录结构

- `main.go`：服务启动入口
- `go.mod`/`go.sum`：Go 依赖管理文件
- `config-debug.yaml`/`config-pro.yaml`：服务配置文件（开发/生产环境）
- `config/`：配置结构体定义
- `global/`：全局变量与配置
- `initialize/`：服务初始化相关（如数据库、日志、配置等）
- `model/`：数据模型定义
- `handler/`：gRPC 业务处理器
- `proto/`：protobuf 协议文件及自动生成代码
- `utils/`：工具类
- `test/`：测试代码
- `temp/`：临时文件夹（如 nacos 配置等）

## 启动方式

1. 配置好 `config-debug.yaml` 或 `config-pro.yaml`。
2. 执行 `go run main.go` 或编译后运行 `main.exe`。
3. 支持通过命令行参数 `-ip` 和 `-port` 指定服务监听地址和端口。

## 主要功能

- 用户信息管理
- gRPC 服务注册与健康检查
- Consul 服务注册与注销
- 日志、数据库等基础设施初始化
 