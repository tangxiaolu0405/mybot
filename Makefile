.PHONY: build build-server build-client clean test install

# 构建目标
BINARY_SERVER=cata
BINARY_CLIENT=catacli
CMD_SERVER=./cmd/cata
CMD_CLIENT=./cmd/catacli

# 默认目标
all: build

# 构建所有
build: build-server build-client

# 构建服务器
build-server:
	@echo "Building $(BINARY_SERVER)..."
	@go build -o $(BINARY_SERVER) $(CMD_SERVER)
	@echo "✓ $(BINARY_SERVER) built successfully"

# 构建客户端
build-client:
	@echo "Building $(BINARY_CLIENT)..."
	@go build -o $(BINARY_CLIENT) $(CMD_CLIENT)
	@echo "✓ $(BINARY_CLIENT) built successfully"

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_SERVER) $(BINARY_CLIENT)
	@echo "✓ Clean completed"

# 运行测试
test:
	@echo "Running tests..."
	@go test ./...

# 安装到系统路径
install: build
	@echo "Installing to $(GOPATH)/bin..."
	@go install $(CMD_SERVER)
	@go install $(CMD_CLIENT)
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
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_SERVER)-linux $(CMD_SERVER)
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_CLIENT)-linux $(CMD_CLIENT)
	@echo "✓ Linux binaries built"

# 交叉编译 - macOS
build-darwin:
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_SERVER)-darwin $(CMD_SERVER)
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_CLIENT)-darwin $(CMD_CLIENT)
	@echo "✓ macOS binaries built"

# 交叉编译 - Windows
build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_SERVER).exe $(CMD_SERVER)
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_CLIENT).exe $(CMD_CLIENT)
	@echo "✓ Windows binaries built"

# 交叉编译 - 所有平台
build-all: build-linux build-darwin build-windows

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  build          - Build both server and client"
	@echo "  build-server   - Build server only (cata)"
	@echo "  build-client   - Build client only (catacli)"
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
