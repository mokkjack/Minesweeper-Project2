
BINARY = program

.PHONY: all build run clean

all: build

build:
	@echo "==> Building $(BINARY)..."
	go build -o $(BINARY) -buildvcs=false

run: build
	@echo "==> Running $(BINARY)..."
	./$(BINARY)

clean:
	@echo "==> Cleaning up..."
	rm -f $(BINARY)