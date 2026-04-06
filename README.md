# gogo: A Go Showcase for Java Developers

Welcome to `gogo`, a single-module Go 1.23 project designed to help Java developers transition to Go. This project showcases Go's core features through interactive demos, with explicit comparisons to Java concepts.

**Key Feature:** Every file in this repository is heavily documented with `// For a Java developer:` comments, explaining Go constructs using their Java equivalents (e.g., Goroutines vs. Threads, Slices vs. ArrayList, etc.).

## 🚀 Getting Started

The easiest way to explore this project is by using the `Makefile`. If you're coming from Maven or Gradle, think of the `Makefile` as your entry point for lifecycle tasks.

### Quick Start
```bash
# Build and run the full showcase
make run
```

### Common Commands (Go vs Maven)

| Task | Go Command | Makefile Target | Java/Maven Equivalent |
|---|---|---|---|
| **Build** | `go build ./cmd/gogo` | `make build` | `mvn package` |
| **Test** | `go test ./...` | `make test` | `mvn test` |
| **Run** | `go run ./cmd/gogo` | `make run` | `mvn spring-boot:run` |
| **Clean** | `rm gogo_binary` | `make clean` | `mvn clean` |
| **Format** | `go fmt ./...` | `make fmt` | `mvn fmt:format` |
| **Check** | `go vet ./...` | `make vet` | `mvn checkstyle:check` |

## 🏗️ Project Architecture

The project is structured to be discoverable:
- `cmd/gogo/main.go`: The main orchestrator that runs all demos in sequence.
- `pkg/`: Contains sub-packages, each focusing on a specific Go concept.

### 📦 Key Demos
- **Concurrency**: Goroutines and Channels (≈ Threads and BlockingQueues).
- **Interfaces**: Implicit implementation (≈ Duck Typing).
- **Generics**: Type-safe collections and helpers.
- **Errors**: Error-as-value handling (≈ Checked Exceptions).
- **Collections**: Functional slice operations (≈ Java Stream API).
- **Patterns**: GoF patterns implemented idiomatically in Go.
- **Advanced**: Reflection, Struct Tags (≈ Annotations), and Build System deep-dive.
- **Common Libraries**: Popular third-party libraries (zap, uuid, gin, testify).

## 🛠️ Go vs Java: Build System & Deployment

### 1. Dependency Management
- **Java**: Uses `pom.xml` (Maven) or `build.gradle`. Dependencies are often complex XML/Groovy/Kotlin trees.
- **Go**: Uses `go.mod`. It's a simple, human-readable text file.
- **How it works**: Run `go mod tidy` to sync dependencies. Go downloads them directly from version control (e.g., GitHub) and caches them in your `$GOPATH/pkg/mod` (similar to `~/.m2/repository`).

### 2. Plugins & Tooling
- **Java**: Heavily reliant on Maven/Gradle plugins for everything (compiling, testing, formatting).
- **Go**: The `go` tool is a "Swiss Army Knife". It includes `fmt` (formatter), `vet` (linter), `test` (test runner), and `doc` (documentation generator) out of the box. No plugin configuration needed.

### 3. Deployment
- **Java**: You typically deploy a `.jar` or `.war` file, which requires a JVM (JRE) and sometimes an Application Server (Tomcat, Wildfly) on the target machine.
- **Go**: Compiles everything into a **single, statically-linked binary**. This binary has NO external dependencies. You just copy it to the server and run it. No "JRE" or "Go Runtime" is needed on the production server.

### 4. Cross-Compilation
Go makes it trivial to build for other operating systems from your local machine:
```bash
# Build for Linux from Mac/Windows
GOOS=linux GOARCH=amd64 go build -o app_linux ./cmd/gogo
```

## 🧪 Running Tests

Go has a built-in testing framework.
- Run all tests: `go test -v ./...`
- Run a specific test: `go test -v ./pkg/generics -run TestMapValues`
- Run benchmarks: `go test -bench=. ./pkg/...`

---
*Happy Gophering! If you're stuck, look for the `// For a Java developer:` comments in the source code.*
