BINARY=cw-agent
OS := $(shell uname -s)
SERVICE_PATH=/etc/systemd/system/$(BINARY).service
PLIST_PATH=~/Library/LaunchAgents/com.xconfig.cw-agent.plist
BINARY_PATH=/usr/local/bin/$(BINARY)

ifeq ($(OS),Linux)
	BINARY_PATH=/usr/local/bin/$(BINARY)
endif

ifeq ($(OS),Darwin)
	BINARY_PATH=/opt/homebrew/bin/$(BINARY)
endif

ifeq ($(OS),Windows_NT)
	BINARY_PATH=$(BINARY).exe
endif

all: build

init:
	@echo "🔍 Checking for Rust toolchain..."
	@if ! command -v cargo >/dev/null 2>&1; then \
		echo "🚀 Rust not found. Installing via rustup..."; \
		curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y; \
		. $$HOME/.cargo/env; \
	fi

	@echo "🛠️ Initializing Rust project..."
	@test -f Cargo.toml || cargo init --vcs none

	@echo "📦 Ensuring required dependencies..."
	@cargo add tokio --features full
	@cargo add serde --features derive
	@cargo add serde_json
	@cargo add reqwest --features json
	@cargo add git2
	@cargo add chrono
	@cargo add clap --features derive
	@cargo add log
	@cargo add thiserror
	@cargo add anyhow

	@echo "📁 Creating module files..."
	@mkdir -p src
	@touch src/{main.rs,scheduler.rs,config.rs,executor.rs,result_store.rs,models.rs}

	@echo "✅ cw-agent init complete."

build:
	cargo build --release

run:
	cargo run -- oneshot

install:
	@echo "📥 Installing binary for $(OS) → $(BINARY_PATH)"
ifeq ($(OS),Windows_NT)
	@echo "⚠️ Windows detected. Please copy manually to a directory in your PATH:"
	@echo "    copy target\\release\\$(BINARY).exe C:\\Program Files\\cw-agent\\"
else
	@install -m 755 target/release/$(BINARY) $(BINARY_PATH)
	@echo "✅ Installed to $(BINARY_PATH)"
endif

install-service: install
	@echo "📦 Installing $(BINARY) as a service..."

	@if [ "$(OS)" = "Linux" ]; then \
		echo "🔧 Installing systemd service on Linux..."; \
		sudo mkdir -p /var/lib/cw-agent; \
		sudo cp example/$(BINARY).service $(SERVICE_PATH); \
		sudo systemctl daemon-reexec; \
		sudo systemctl enable $(BINARY); \
		sudo systemctl restart $(BINARY); \
		systemctl status $(BINARY); \
	elif [ "$(OS)" = "Darwin" ]; then \
		echo "🔧 Installing launchd agent on macOS..."; \
		mkdir -p /var/lib/cw-agent; \
		cp example/$(BINARY).plist $(PLIST_PATH); \
		launchctl unload $(PLIST_PATH) 2>/dev/null || true; \
		launchctl load $(PLIST_PATH); \
		launchctl start com.xconfig.cw-agent; \
		launchctl list | grep com.xconfig.cw-agent; \
	elif [ "$(OS)" = "Windows_NT" ]; then \
		echo "⚠️ Windows install-service is not supported via Makefile directly."; \
		echo "👉 Use 'sc create' or NSSM to install manually."; \
	else \
		echo "❌ Unsupported OS: $(OS)"; \
		exit 1; \
	fi

	@echo "✅ install-service complete."

uninstall-service:
	@echo "🧹 Uninstalling service..."
	@if [ "$(OS)" = "Linux" ]; then \
		sudo systemctl stop $(BINARY); \
		sudo systemctl disable $(BINARY); \
		sudo rm -f $(SERVICE_PATH); \
		sudo systemctl daemon-reexec; \
		echo "✅ Linux systemd service uninstalled."; \
	elif [ "$(OS)" = "Darwin" ]; then \
		launchctl unload $(PLIST_PATH) 2>/dev/null || true; \
		rm -f $(PLIST_PATH); \
		echo "✅ macOS launchd agent uninstalled."; \
	else \
		echo "⚠️ Manual uninstall required on $(OS)."; \
	fi

start-service:
	@echo "🚀 Starting $(BINARY) service..."
	@if [ "$(OS)" = "Linux" ]; then sudo systemctl start $(BINARY); fi
	@if [ "$(OS)" = "Darwin" ]; then launchctl start com.xconfig.cw-agent; fi

stop-service:
	@echo "🛑 Stopping $(BINARY) service..."
	@if [ "$(OS)" = "Linux" ]; then sudo systemctl stop $(BINARY); fi
	@if [ "$(OS)" = "Darwin" ]; then launchctl stop com.xconfig.cw-agent; fi

restart-service:
	@echo "♻️ Restarting $(BINARY) service..."
	@if [ "$(OS)" = "Linux" ]; then sudo systemctl restart $(BINARY); fi
	@if [ "$(OS)" = "Darwin" ]; then launchctl stop com.xconfig.cw-agent && launchctl start com.xconfig.cw-agent; fi

logs:
	@echo "📜 Viewing logs..."
	@if [ "$(OS)" = "Linux" ]; then journalctl -u $(BINARY) -f; fi
	@if [ "$(OS)" = "Darwin" ]; then tail -f /tmp/cw-agent.out /tmp/cw-agent.err; fi

test:
	cargo test

fmt:
	cargo fmt

lint:
	cargo clippy -- -D warnings

clean:
	cargo clean
