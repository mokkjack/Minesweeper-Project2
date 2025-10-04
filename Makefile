# Makefile for Minesweeper-Project2 by Zhang

#This makefile will actually compile the program in Linux

BINARY = program
OS := $(shell uname -s 2>/dev/null || echo Windows)

.PHONY: all deps sysdeps godeps build run clean

all: build

# System dependency check (Linux only)
sysdeps:
ifeq ($(OS),Linux)
	@echo "==> Checking system dependencies..."
	@missing=""; \
	for pkg in libgl1-mesa-dev libglu1-mesa-dev xorg-dev pkg-config; do \
		if ! dpkg -s $$pkg >/dev/null 2>&1; then \
			missing="$$missing $$pkg"; \
		fi; \
	done; \
	if [ -n "$$missing" ]; then \
		echo "Missing system packages:$$missing"; \
		echo "Updating package lists..."; \
		sudo apt-get update -y; \
		echo "Installing missing dependencies"; \
		sudo apt-get install -y $$missing; \
	else \
		echo "All system dependencies are installed."; \
	fi
else
	@echo "Please Use Linux or Linux subsystem for Window, this program will not work in Window!"
endif

# Go dependencies
godeps:
	@echo "==> Checking Go module dependencies..."
	@if [ ! -f go.mod ]; then \
		echo "No go.mod found. Initializing Go module..."; \
		go mod init minesweeper || true; \
	fi
	go mod tidy

deps: sysdeps godeps
	@echo "==> All dependencies satisfied."

build: deps
	@echo "==> Building $(BINARY)..."
	go build -o $(BINARY) -buildvcs=false

run: build
	@echo "==> Running $(BINARY)..."
	./$(BINARY)

clean:
	@echo "==> Cleaning up..."
	rm -f $(BINARY)
