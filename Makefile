.PHONY: gqlgen
gqlgen:
	gqlgen --verbose

.PHONY: build
build: go-build

.PHONY: clean
clean: go-clean ## Clean build cache and dependencies

.PHONY: wire
wire:
	@cd cmd/graphy/inject && wire

go-build:
	@echo "Building Go services..."
	@rm -rf build
	@mkdir build
	go build -o build -v ./...
	@echo "Go services available at ./build"

go-clean: go-clean-cache go-clean-deps

go-clean-cache:
	@echo "Cleaning build cache..."
	go clean -cache

go-clean-test-cache:
	@echo "Cleaning test cache..."
	go clean -testcache

go-clean-deps:
	@echo "Cleaning dependencies..."
	go mod tidy

go-deps:
	@echo "Installing dependencies..."
	go mod download
