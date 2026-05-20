.PHONY: build clean test install

BINARY=cata
CMD=./cmd/cata

all: build

build:
	@echo "Building $(BINARY)..."
	@go build -o $(BINARY) $(CMD)
	@echo "✓ $(BINARY) built successfully"

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY) catacli catacli.exe
	@echo "✓ Clean completed"

# 运行测试
test:
	@echo "Running tests..."
	@go test ./...

# 安装到系统路径
install: build
	@echo "Installing to $(GOPATH)/bin..."
	@go install $(CMD)
	@echo "✓ Installed successfully"

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format completed"

# 检查代码
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet completed"

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies updated"

# 交叉编译 - Linux
build-linux:
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux $(CMD)
	@echo "✓ Linux binary built"

# 交叉编译 - macOS
build-darwin:
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin $(CMD)
	@echo "✓ macOS binary built"

# 交叉编译 - Windows
build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY).exe $(CMD)
	@echo "✓ Windows binary built"

# 交叉编译 - 所有平台
build-all: build-linux build-darwin build-windows

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  build          - Build cata binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  deps           - Download and update dependencies"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-darwin   - Build for macOS"
	@echo "  build-windows  - Build for Windows"
	@echo "  build-all      - Build for all platforms"
