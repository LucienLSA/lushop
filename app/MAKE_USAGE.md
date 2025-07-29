# LuShop Makefile 使用指南 (Linux版本)

## 快速开始

### 1. 编译所有服务
```bash
# 在app目录下执行
make all
```

### 2. 编译单个服务
```bash
# 编译API网关服务
make build-api

# 编译商品服务
make build-goods

# 编译指定服务
make build-service SERVICE=goods_srv
```

### 3. 查看帮助
```bash
make help
```

## 详细命令说明

### 编译相关命令

| 命令 | 说明 |
|------|------|
| `make all` | 清理并编译所有服务 |
| `make build-all` | 编译所有服务 |
| `make build-api` | 编译API网关服务 |
| `make build-services` | 编译所有微服务 |
| `make build-service SERVICE=服务名` | 编译指定微服务 |
| `make build-goods` | 编译商品服务 |
| `make build-inventory` | 编译库存服务 |
| `make build-order` | 编译订单服务 |
| `make build-user` | 编译用户服务 |
| `make build-userop` | 编译用户操作服务 |

### 依赖管理

| 命令 | 说明 |
|------|------|
| `make deps` | 下载所有依赖 |
| `make deps-update` | 更新所有依赖 |

### 代码质量

| 命令 | 说明 |
|------|------|
| `make test` | 运行所有测试 |
| `make fmt` | 格式化代码 |
| `make vet` | 代码检查 |

### 系统管理

| 命令 | 说明 |
|------|------|
| `make clean` | 清理构建文件 |
| `make install` | 安装到系统路径 |
| `make uninstall` | 卸载服务 |
| `make cross-build` | 交叉编译 |

### 信息查看

| 命令 | 说明 |
|------|------|
| `make help` | 显示帮助信息 |
| `make info` | 显示构建信息 |
| `make status` | 显示服务状态 |
| `make check-go` | 检查Go环境 |

## 构建输出

编译后的二进制文件将保存在 `build/bin/` 目录下：

```
build/bin/
├── lushop-api      # API网关服务
├── goods_srv       # 商品服务
├── inventory_srv   # 库存服务
├── order_srv       # 订单服务
├── user_srv        # 用户服务
└── userop_srv      # 用户操作服务
```

## 环境变量

可以通过环境变量自定义构建参数：

```bash
# 设置版本号
export VERSION=2.0.0

# 设置目标平台
export GOOS=linux
export GOARCH=amd64

# 编译
make all
```

## 常见问题

### 1. 权限问题
如果遇到权限问题，确保有执行权限：
```bash
chmod +x build/bin/*
```

### 2. 依赖问题
如果编译失败，先更新依赖：
```bash
make deps-update
make all
```

### 3. 清理问题
如果构建有问题，清理后重新编译：
```bash
make clean
make all
```

### 4. 系统安装
安装到系统路径（需要sudo权限）：
```bash
make install
```

## 高级用法

### 1. 交叉编译
为不同平台编译：
```bash
make cross-build
```

### 2. 自定义服务编译
编译特定服务：
```bash
make build-service SERVICE=goods_srv
```

### 3. 检查服务状态
查看已编译的服务：
```bash
make status
```

## 注意事项

1. **Go版本要求**: 确保Go版本 >= 1.24
2. **权限要求**: 某些命令需要sudo权限
3. **网络要求**: 下载依赖需要网络连接
4. **磁盘空间**: 确保有足够的磁盘空间存储编译文件

## 故障排除

### 编译失败
```bash
# 检查Go环境
make check-go

# 更新依赖
make deps-update

# 清理重新编译
make clean
make all
```

### 依赖下载失败
```bash
# 设置Go代理（如果需要）
export GOPROXY=https://goproxy.cn,direct

# 重新下载依赖
make deps
```

### 权限问题
```bash
# 确保有执行权限
chmod +x build/bin/*

# 或者使用sudo
sudo make install
``` 