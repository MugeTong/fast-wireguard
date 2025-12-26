VERSION := 0.1.0
FULL_NAME := "Fast WireGuard"


build:
	@mkdir -p releases && \
	makeself ./src ./releases/fast-wireguard-$(VERSION)-Linux-x86_64.sh $(FULL_NAME) ./setup.sh

all:
	@mkdir -p releases && \
	makeself ./src ./releases/fast-wireguard-$(VERSION)-Linux-x86_64.sh $(FULL_NAME) ./setup.sh && \
	./releases/fast-wireguard-$(VERSION)-Linux-x86_64.sh
