# Makefile for worldclocks

# Variables
BINARY_NAME=worldclocks
INSTALL_PATH=/usr/local/bin
GO=go
GOFLAGS=-trimpath
LDFLAGS=-s -w

# Default target
.DEFAULT_GOAL := build

# Build the binary
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	${GO} build ${GOFLAGS} -ldflags="${LDFLAGS}" -o ${BINARY_NAME} .
	@echo "Build complete: ${BINARY_NAME}"

# Build for all platforms (useful for testing before release)
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 ${GO} build ${GOFLAGS} -ldflags="${LDFLAGS}" -o ${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 ${GO} build ${GOFLAGS} -ldflags="${LDFLAGS}" -o ${BINARY_NAME}-linux-arm64 .
	GOOS=darwin GOARCH=amd64 ${GO} build ${GOFLAGS} -ldflags="${LDFLAGS}" -o ${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 ${GO} build ${GOFLAGS} -ldflags="${LDFLAGS}" -o ${BINARY_NAME}-darwin-arm64 .
	GOOS=windows GOARCH=amd64 ${GO} build ${GOFLAGS} -ldflags="${LDFLAGS}" -o ${BINARY_NAME}-windows-amd64.exe .
	@echo "Cross-platform build complete"

# Install the binary to system path (requires sudo)
# Note: Run 'make build' first to create the binary
.PHONY: install
install:
	@if [ ! -f ${BINARY_NAME} ]; then \
		echo "Error: ${BINARY_NAME} binary not found. Run 'make build' first."; \
		exit 1; \
	fi
	@echo "Installing ${BINARY_NAME} to ${INSTALL_PATH}..."
	@install -d ${INSTALL_PATH}
	@install -m 755 ${BINARY_NAME} ${INSTALL_PATH}/${BINARY_NAME}
	@echo "Installation complete. Run '${BINARY_NAME}' to start."

# Uninstall the binary from system path (requires sudo)
.PHONY: uninstall
uninstall:
	@echo "Uninstalling ${BINARY_NAME} from ${INSTALL_PATH}..."
	@rm -f ${INSTALL_PATH}/${BINARY_NAME}
	@echo "Uninstall complete"

# Run the application
.PHONY: run
run:
	@${GO} run .

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f ${BINARY_NAME}
	@rm -f ${BINARY_NAME}-*
	@rm -rf dist/
	@echo "Clean complete"

# Run tests (when tests are added)
.PHONY: test
test:
	@${GO} test -v ./...

# Update dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@${GO} mod download
	@${GO} mod tidy
	@echo "Dependencies updated"

# Test GoReleaser configuration without publishing
.PHONY: release-test
release-test:
	@echo "Testing GoReleaser configuration..."
	@goreleaser release --snapshot --clean
	@echo "Release test complete. Check dist/ directory."

# Create a release (requires goreleaser and git tag)
.PHONY: release
release:
	@echo "Creating release with GoReleaser..."
	@goreleaser release --clean
	@echo "Release complete"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-all    - Build for all supported platforms"
	@echo "  install      - Install binary to ${INSTALL_PATH} (requires sudo)"
	@echo "  uninstall    - Remove binary from ${INSTALL_PATH} (requires sudo)"
	@echo "  run          - Run the application without building"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Update Go dependencies"
	@echo "  release-test - Test GoReleaser configuration locally"
	@echo "  release      - Create a release with GoReleaser"
	@echo "  help         - Show this help message"
