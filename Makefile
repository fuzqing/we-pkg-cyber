# we-pkg-cyber/Makefile

.PHONY: build-wasm dev clean install

# 定义输出路径
WASM_OUT = web/public/main.wasm

# 【关键修复】使用 $(shell ...) 让 Make 去执行终端命令，而不是读取变量
GOROOT_PATH = $(shell go env GOROOT)
WASM_EXEC_SRC = $(GOROOT_PATH)/lib/wasm/wasm_exec.js
WASM_EXEC_DEST = web/public/wasm_exec.js

build-wasm:
	@echo "=> 🔨 正在编译 Golang Wasm 引擎..."
	cd wasm-core && GOOS=js GOARCH=wasm go build -o ../$(WASM_OUT) cmd/wepkg/main.go
	@echo "=> 🚚 正在自动注入最新版 wasm_exec.js..."
	cp "$(WASM_EXEC_SRC)" "$(WASM_EXEC_DEST)"
	@echo "=> ✅ 底层引擎构建与注入完成！"

dev: build-wasm
	@echo "=> 🚀 启动前端赛博控制台..."
	cd web && npm run dev

install:
	@echo "=> 📦 安装前端依赖..."
	cd web && npm install
	
clean:
	@echo "=> 🧹 清理构建产物..."
	rm -f $(WASM_OUT) $(WASM_EXEC_DEST)