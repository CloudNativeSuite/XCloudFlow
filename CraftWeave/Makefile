APP_NAME := craftweave
MAIN_FILE := main.go
MODULE := craftweave
TAG ?= latest
PLAYBOOK ?= examples/deploy.yaml

.PHONY: all build run clean init ansible playbook vault cmdb plugin help

all: build

init:
	go mod init craftweave
	GOPROXY=https://goproxy.cn,direct go get github.com/spf13/cobra@latest
	GOPROXY=https://goproxy.cn,direct go get golang.org/x/crypto/ssh
	GOPROXY=https://goproxy.cn,direct go get gopkg.in/yaml.v3
	go mod tidy

build:
	go build -o $(APP_NAME) $(MAIN_FILE)

run:
	go run $(MAIN_FILE)

ansible:
	go run $(MAIN_FILE) ansible all -m ping

playbook:
	go run $(MAIN_FILE) playbook $(PLAYBOOK)

vault:
	go run $(MAIN_FILE) vault

cmdb:
	go run $(MAIN_FILE) cmdb export --format json > cmdb-graph.json

plugin:
	go run $(MAIN_FILE) plugin run ./plugins/hello.wasm

clean:
	rm -f $(APP_NAME)

help:
	@echo "🧶 Xconfig CLI Usage"
	@echo ""
	@echo "make build                编译 craftweave 可执行文件"
	@echo "make run                  运行默认入口"
	@echo "make ansible              执行 ansible all -m ping"
	@echo "make playbook             执行默认 playbook 文件"
	@echo "make vault                加解密命令（占位）"
	@echo "make cmdb                 输出图数据库格式"
	@echo "make plugin               加载运行插件（WASM）"
	@echo "make init                 初始化依赖（go mod tidy）"
	@echo "make clean                清理构建产物"
