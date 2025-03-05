# Makefile for compiling .proto files into Go code recursively

# Directories
PROTO_DIR := proto
GO_OUT := types

# Recursively find all .proto files in the PROTO_DIR
PROTO_FILES := $(shell find $(PROTO_DIR) -type f -name "*.proto")

# Check if tools are installed
PROTOC := $(shell command -v protoc 2>/dev/null)
PROTOC_GEN_GO := $(shell command -v protoc-gen-go 2>/dev/null)
PROTOC_GEN_GO_GRPC := $(shell command -v protoc-gen-go-grpc 2>/dev/null)

.PHONY: all check_tools proto clean

# Default target depends on checking tools and compiling protos
all: check_tools proto

# Check for required tools and install missing Go plugins automatically.
check_tools:
ifeq ($(PROTOC),)
	$(error "protoc is not installed. Please install the Protocol Buffers compiler.")
endif
ifeq ($(PROTOC_GEN_GO),)
	@echo "protoc-gen-go not found. Installing..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
endif
ifeq ($(PROTOC_GEN_GO_GRPC),)
	@echo "protoc-gen-go-grpc not found. Installing..."
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
endif

# Compile all .proto files into Go source files
proto: check_tools
	@mkdir -p $(GO_OUT)
	$(PROTOC) --proto_path=$(PROTO_DIR) \
	          --go_out=paths=source_relative:$(GO_OUT) \
	          --go-grpc_out=paths=source_relative:$(GO_OUT) \
	          $(PROTO_FILES)

# Clean generated files
clean:
	rm -rf $(GO_OUT)