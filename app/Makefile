# LuShop 微服务编译 Makefile (Linux版本)
# 用于编译 lushop_api 和 lushop_srvs 中的所有 main.go 文件

# 变量定义
GO=go
GOOS?=linux
GOARCH?=amd64
CGO_ENABLED=0

# 构建目录
BUILD_DIR=build
BIN_DIR=$(BUILD_DIR)/bin

# 服务列表
API_SERVICE=lushop_api
SERVICES=goods_srv inventory_srv order_srv user_srv userop_srv

# 版本信息
VERSION?=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 编译参数
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -s -w"

# 默认目标
.PHONY: all
all: clean build-all

# 清理构建目录
.PHONY: clean
clean:
	@echo "清理构建目录..."
	@rm -rf $(BUILD_DIR)
	@rm -f lushop_api/main
	@for service in $(SERVICES); do \
		rm -f lushop_srvs/$$service/main; \
	done
	@echo "清理完成"

# 创建构建目录
.PHONY: mkdir
mkdir:
	@mkdir -p $(BIN_DIR)

# 编译所有服务
.PHONY: build-all
build-all: mkdir build-api build-services
	@echo "所有服务编译完成"

# 编译API网关服务
.PHONY: build-api
build-api:
	@echo "编译 API 网关服务..."
	@cd lushop_api && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_api/main $(BIN_DIR)/lushop-api
	@chmod +x $(BIN_DIR)/lushop-api
	@echo "API 网关服务编译完成: $(BIN_DIR)/lushop-api"

# 编译所有微服务
.PHONY: build-services
build-services:
	@echo "编译微服务..."
	@for service in $(SERVICES); do \
		echo "编译 $$service 服务..."; \
		cd lushop_srvs/$$service && $(GO) build $(LDFLAGS) -o main .; \
		cp main ../../$(BIN_DIR)/$$service; \
		chmod +x ../../$(BIN_DIR)/$$service; \
		echo "$$service 服务编译完成: $(BIN_DIR)/$$service"; \
	done

# 编译单个API服务
.PHONY: build-api-only
build-api-only: mkdir
	@echo "编译 API 网关服务..."
	@cd lushop_api && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_api/main $(BIN_DIR)/lushop-api
	@chmod +x $(BIN_DIR)/lushop-api
	@echo "API 网关服务编译完成: $(BIN_DIR)/lushop-api"

# 编译单个微服务
.PHONY: build-service
build-service: mkdir
	@if [ -z "$(SERVICE)" ]; then \
		echo "错误: 请指定要编译的服务名称"; \
		echo "使用方法: make build-service SERVICE=goods_srv"; \
		exit 1; \
	fi
	@echo "编译 $(SERVICE) 服务..."
	@cd lushop_srvs/$(SERVICE) && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_srvs/$(SERVICE)/main $(BIN_DIR)/$(SERVICE)
	@chmod +x $(BIN_DIR)/$(SERVICE)
	@echo "$(SERVICE) 服务编译完成: $(BIN_DIR)/$(SERVICE)"

# 编译商品服务
.PHONY: build-goods
build-goods: mkdir
	@echo "编译商品服务..."
	@cd lushop_srvs/goods_srv && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_srvs/goods_srv/main $(BIN_DIR)/goods_srv
	@chmod +x $(BIN_DIR)/goods_srv
	@echo "商品服务编译完成: $(BIN_DIR)/goods_srv"

# 编译库存服务
.PHONY: build-inventory
build-inventory: mkdir
	@echo "编译库存服务..."
	@cd lushop_srvs/inventory_srv && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_srvs/inventory_srv/main $(BIN_DIR)/inventory_srv
	@chmod +x $(BIN_DIR)/inventory_srv
	@echo "库存服务编译完成: $(BIN_DIR)/inventory_srv"

# 编译订单服务
.PHONY: build-order
build-order: mkdir
	@echo "编译订单服务..."
	@cd lushop_srvs/order_srv && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_srvs/order_srv/main $(BIN_DIR)/order_srv
	@chmod +x $(BIN_DIR)/order_srv
	@echo "订单服务编译完成: $(BIN_DIR)/order_srv"

# 编译用户服务
.PHONY: build-user
build-user: mkdir
	@echo "编译用户服务..."
	@cd lushop_srvs/user_srv && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_srvs/user_srv/main $(BIN_DIR)/user_srv
	@chmod +x $(BIN_DIR)/user_srv
	@echo "用户服务编译完成: $(BIN_DIR)/user_srv"

# 编译用户操作服务
.PHONY: build-userop
build-userop: mkdir
	@echo "编译用户操作服务..."
	@cd lushop_srvs/userop_srv && $(GO) build $(LDFLAGS) -o main .
	@cp lushop_srvs/userop_srv/main $(BIN_DIR)/userop_srv
	@chmod +x $(BIN_DIR)/userop_srv
	@echo "用户操作服务编译完成: $(BIN_DIR)/userop_srv"

# 下载依赖
.PHONY: deps
deps:
	@echo "下载 API 网关服务依赖..."
	@cd lushop_api && $(GO) mod download
	@echo "下载微服务依赖..."
	@for service in $(SERVICES); do \
		echo "下载 $$service 依赖..."; \
		cd lushop_srvs/$$service && $(GO) mod download; \
	done
	@echo "所有依赖下载完成"

# 更新依赖
.PHONY: deps-update
deps-update:
	@echo "更新 API 网关服务依赖..."
	@cd lushop_api && $(GO) mod tidy
	@echo "更新微服务依赖..."
	@for service in $(SERVICES); do \
		echo "更新 $$service 依赖..."; \
		cd lushop_srvs/$$service && $(GO) mod tidy; \
	done
	@echo "所有依赖更新完成"

# 测试
.PHONY: test
test:
	@echo "运行 API 网关服务测试..."
	@cd lushop_api && $(GO) test ./...
	@echo "运行微服务测试..."
	@for service in $(SERVICES); do \
		echo "运行 $$service 测试..."; \
		cd lushop_srvs/$$service && $(GO) test ./...; \
	done
	@echo "所有测试完成"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化 API 网关服务代码..."
	@cd lushop_api && $(GO) fmt ./...
	@echo "格式化微服务代码..."
	@for service in $(SERVICES); do \
		echo "格式化 $$service 代码..."; \
		cd lushop_srvs/$$service && $(GO) fmt ./...; \
	done
	@echo "代码格式化完成"

# 代码检查
.PHONY: vet
vet:
	@echo "检查 API 网关服务代码..."
	@cd lushop_api && $(GO) vet ./...
	@echo "检查微服务代码..."
	@for service in $(SERVICES); do \
		echo "检查 $$service 代码..."; \
		cd lushop_srvs/$$service && $(GO) vet ./...; \
	done
	@echo "代码检查完成"

# 交叉编译 (用于不同平台)
.PHONY: cross-build
cross-build: clean
	@echo "交叉编译所有服务..."
	@mkdir -p $(BIN_DIR)
	@echo "编译 API 网关服务 (Linux)..."
	@cd lushop_api && GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o main .
	@cp lushop_api/main $(BIN_DIR)/lushop-api-linux-amd64
	@chmod +x $(BIN_DIR)/lushop-api-linux-amd64
	@echo "编译微服务 (Linux)..."
	@for service in $(SERVICES); do \
		echo "编译 $$service 服务 (Linux)..."; \
		cd lushop_srvs/$$service && GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o main .; \
		cp main ../../$(BIN_DIR)/$$service-linux-amd64; \
		chmod +x ../../$(BIN_DIR)/$$service-linux-amd64; \
	done
	@echo "交叉编译完成"

# 安装到系统路径
.PHONY: install
install: build-all
	@echo "安装服务到系统路径..."
	@sudo cp $(BIN_DIR)/* /usr/local/bin/
	@echo "安装完成"

# 卸载
.PHONY: uninstall
uninstall:
	@echo "卸载服务..."
	@sudo rm -f /usr/local/bin/lushop-api
	@for service in $(SERVICES); do \
		sudo rm -f /usr/local/bin/$$service; \
	done
	@echo "卸载完成"

# 显示帮助信息
.PHONY: help
help:
	@echo "LuShop 微服务编译 Makefile (Linux版本)"
	@echo ""
	@echo "可用目标:"
	@echo "  all              - 编译所有服务"
	@echo "  clean            - 清理构建文件"
	@echo "  build-all        - 编译所有服务"
	@echo "  build-api        - 编译API网关服务"
	@echo "  build-services   - 编译所有微服务"
	@echo "  build-service    - 编译指定微服务 (需要指定 SERVICE=服务名)"
	@echo "  build-goods      - 编译商品服务"
	@echo "  build-inventory  - 编译库存服务"
	@echo "  build-order      - 编译订单服务"
	@echo "  build-user       - 编译用户服务"
	@echo "  build-userop     - 编译用户操作服务"
	@echo "  deps             - 下载所有依赖"
	@echo "  deps-update      - 更新所有依赖"
	@echo "  test             - 运行所有测试"
	@echo "  fmt              - 格式化代码"
	@echo "  vet              - 代码检查"
	@echo "  cross-build      - 交叉编译 (Linux)"
	@echo "  install          - 安装到系统路径"
	@echo "  uninstall        - 卸载服务"
	@echo "  help             - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make all                    # 编译所有服务"
	@echo "  make build-service SERVICE=goods_srv  # 编译商品服务"
	@echo "  make clean                  # 清理构建文件"
	@echo "  make deps                   # 下载依赖"
	@echo "  make install                # 安装到系统"

# 显示构建信息
.PHONY: info
info:
	@echo "构建信息:"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  Git提交: $(GIT_COMMIT)"
	@echo "  目标平台: $(GOOS)/$(GOARCH)"
	@echo "  构建目录: $(BUILD_DIR)"
	@echo "  二进制目录: $(BIN_DIR)"

# 检查Go环境
.PHONY: check-go
check-go:
	@echo "检查Go环境..."
	@$(GO) version
	@echo "Go环境检查完成"

# 显示服务状态
.PHONY: status
status:
	@echo "服务状态:"
	@echo "API网关服务:"
	@if [ -f "$(BIN_DIR)/lushop-api" ]; then \
		echo "  ✓ lushop-api 已编译"; \
		ls -lh $(BIN_DIR)/lushop-api; \
	else \
		echo "  ✗ lushop-api 未编译"; \
	fi
	@echo "微服务:"
	@for service in $(SERVICES); do \
		if [ -f "$(BIN_DIR)/$$service" ]; then \
			echo "  ✓ $$service 已编译"; \
			ls -lh $(BIN_DIR)/$$service; \
		else \
			echo "  ✗ $$service 未编译"; \
		fi; \
	done 