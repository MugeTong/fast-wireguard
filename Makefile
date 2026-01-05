VERSION := 0.2.0
BINARY_NAME := "fwg"
PACKAGE_NAME := "fast-wireguard"
DISPLAY_NAME := "Fast-WireGuard"
BUILD_DIR := "./releases/build"
RELEASE_DIR := "./releases"

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))


build:
	@mkdir -p $(BUILD_DIR) && \
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fwg/main.go

pack:
	@mkdir -p $(RELEASE_DIR) && \
	cp ./setup.sh $(BUILD_DIR)/ && \
	makeself $(BUILD_DIR) $(RELEASE_DIR)/$(PACKAGE_NAME)-Linux-amd64.sh $(DISPLAY_NAME) ./setup.sh

# This command is used for text
run:
	@go run ./cmd/fwg/main.go $(ARGS)


# To prevent make from attempting to build a second target, add the catch-all rule
%:
	@:
