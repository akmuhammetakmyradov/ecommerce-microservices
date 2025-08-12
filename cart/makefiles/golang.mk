APP_NAME := cart
BIN_DIR := bin
PORT := 8080
LOCAL_BIN := $(CURDIR)/$(BIN_DIR)
VENDOR_PROTO_DIR := vendor.protogen

.PHONY: build run test clean install-deps vendor-proto generate_protoc generate_stocks_protoc

build: ## 🔨 Build the cart app
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/main.go

run: build ## 🚀 Run the cart app locally
	@echo "🚀 Running $(APP_NAME) on port $(PORT)..."
	@./$(BIN_DIR)/$(APP_NAME) -port $(PORT)

test: ## 🧪 Run unit tests
	@echo "🧪 Testing $(APP_NAME)..."
	@go test -v ./...

clean: ## 🧹 Clean build files and vendored protos
	@echo "🧹 Cleaning $(APP_NAME) build files and vendored protos..."
	@rm -rf $(BIN_DIR) $(VENDOR_PROTO_DIR)
	@go clean

install-deps: ## 📥 Install protoc-gen-* plugins locally
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

vendor-proto: ## 📦 Vendor google api protos if not already present
	@if [ ! -d $(VENDOR_PROTO_DIR)/google/api ]; then \
		git clone https://github.com/googleapis/googleapis.git $(VENDOR_PROTO_DIR)/googleapis && \
		mkdir -p $(VENDOR_PROTO_DIR)/google && \
		mv $(VENDOR_PROTO_DIR)/googleapis/google/api $(VENDOR_PROTO_DIR)/google/ && \
		rm -rf $(VENDOR_PROTO_DIR)/googleapis; \
	fi

generate_stocks_protoc:
	protoc -I ../proto -I $(VENDOR_PROTO_DIR) \
    --go_out=pkg/api --go_opt=paths=source_relative \
    --go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
    ../proto/stocks/stocks.proto

generate_protoc: install-deps vendor-proto generate_stocks_protoc ## 🛠 Generate protobuf, grpc and grpc-gateway code
	protoc -I ../proto -I $(VENDOR_PROTO_DIR) \
    --go_out=. --go_opt=module=cart \
    --go-grpc_out=. --go-grpc_opt=module=cart \
    --grpc-gateway_out=. --grpc-gateway_opt=module=cart \
		../proto/cart/cart.proto
