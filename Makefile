# Makefile for gogo project
# This provides a unified entry point for CLI tasks, similar to how a Java developer might use Maven or Gradle.

# .PHONY is used to declare that these targets are not files.
# Java equivalent: This is like defining Maven goals or Gradle tasks.
.PHONY: all build test run clean fmt vet lint help

# Default target
# Java equivalent: This is like the 'mvn install' or 'mvn compile' goal.
all: build test

## build: Build the project binary
# Java equivalent: 'mvn package' (creates a single executable binary, like a Fat JAR but with no JVM needed).
build:
	@echo "Building gogo binary..."
	go build -o gogo_binary ./cmd/gogo

## test: Run all tests
# Java equivalent: 'mvn test' (runs all files ending in _test.go).
test:
	@echo "Running all tests..."
	go test -v ./...

## run: Build and run the project
# Java equivalent: 'mvn spring-boot:run' or 'java -jar target/app.jar'.
run: build
	@echo "Running gogo..."
	./gogo_binary

## clean: Remove build artifacts
# Java equivalent: 'mvn clean' (deletes the binary).
clean:
	@echo "Cleaning up..."
	rm -f gogo_binary

## fmt: Format all Go files
# Java equivalent: 'mvn fmt:format' (Go has this built-in).
fmt:
	@echo "Formatting Go source code..."
	go fmt ./...

## vet: Run go vet on all packages
# Java equivalent: A form of static analysis (Checkstyle/PMD).
vet:
	@echo "Vetting Go source code..."
	go vet ./...

## lint: Placeholder for linting (requires golangci-lint)
# Java equivalent: Advanced static analysis (SonarQube/FindBugs).
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run; \
		else \
		echo "golangci-lint not installed. Skipping..."; \
	fi

## help: Display help information
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
