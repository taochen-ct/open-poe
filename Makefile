# Simple Makefile for a Go project

wire:
	@cd ./cmd/api && wire

# Build the application
all: build

build:
	@echo "Building..."

	@go build -o main ./cmd/api

	@echo "Build Completed"


static-build:
	@echo "Static Building..."

	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main ./cmd/api

	@echo "Static Build Completed"

# Run the application
run:
	@go run ./cmd/api

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi


.PHONY: all build static-build run clean watch