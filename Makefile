.PHONY: help all xcloud-init xcloud-build xcloud-run xcloud-up xcloud-down xcloud-export xcloud-import xcloud-ansible xconfig-init xconfig-build xconfig-run xconfig-playbook xconfig-agent-build xconfig-agent-run xconfig-agent-install

ENV ?= sit

help:
	@echo "🚀 Project Targets"
	@echo "  make xcloud-build          # 构建 Go 版 XCloud CLI"
	@echo "  make xcloud-run ENV=sit    # 运行 XCloud CLI (示例)"
	@echo "  make xconfig-build         # 构建 Go 版 Xconfig"
	@echo "  make xconfig-playbook      # 使用默认示例执行 playbook"
	@echo "  make xconfig-agent-build   # 构建 Rust 版 xconfig-agent"
	@echo "  make xconfig-agent-run     # 运行 xconfig-agent oneshot"

all: help

build: xcloud-build xconfig-build xconfig-agent-build

xcloud-init:
	$(MAKE) -C xcloud-cli init

xcloud-build:
	$(MAKE) -C xcloud-cli build

xcloud-run:
	$(MAKE) -C xcloud-cli run ENV=$(ENV)

xcloud-up:
	$(MAKE) -C xcloud-cli up ENV=$(ENV)

xcloud-down:
	$(MAKE) -C xcloud-cli down ENV=$(ENV)

xcloud-export:
	$(MAKE) -C xcloud-cli export ENV=$(ENV)

xcloud-import:
	$(MAKE) -C xcloud-cli import ENV=$(ENV)

xcloud-ansible:
	$(MAKE) -C xcloud-cli ansible ENV=$(ENV)

xconfig-init:
	$(MAKE) -C xconfig init

xconfig-build:
	$(MAKE) -C xconfig build

xconfig-run:
	$(MAKE) -C xconfig run

xconfig-playbook:
	$(MAKE) -C xconfig playbook

xconfig-agent-build:
	$(MAKE) -C xconfig-agent build

xconfig-agent-run:
	$(MAKE) -C xconfig-agent run

xconfig-agent-install:
	$(MAKE) -C xconfig-agent install
