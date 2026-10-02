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
- **Go Syntax**: Core syntax differences for Java developers (Variables, Pointers, multiple returns).
- **Concurrency**: Goroutines and Channels (≈ Threads and BlockingQueues).
- **Interfaces**: Implicit implementation (≈ Duck Typing).
- **Generics**: Type-safe collections and helpers.
- **Errors**: Error-as-value handling (≈ Checked Exceptions).
- **Collections**: Functional slice operations (≈ Java Stream API).
- **Patterns & Antipatterns**: GoF patterns implemented idiomatically, plus 10 crucial Go antipatterns that Java developers must unlearn.
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

## 🧩 Go Patterns & Antipatterns Guide (for Java Developers)

Java developers bring years of object-oriented assumptions, JVM runtime knowledge, and Spring architectural patterns to Go. While these patterns are effective in Java, applying them directly to Go leads to subtle bugs, performance degradation, and unidiomatic code.

### 🌟 Idiomatic Go Patterns

1. **Functional Options Pattern**:
   - *Java Equivalent*: Builder Pattern (`new ServerConfigBuilder().port(8080).build()`) or overloaded constructors.
   - *Go Idiom*: Variadic functions accepting closures `func WithTimeout(d time.Duration) Option`. Allows extensible, readable configuration with sensible defaults without constructor explosion.
2. **Consumer-Driven Interfaces ("Accept Interfaces, Return Structs")**:
   - *Java Equivalent*: Producer interfaces (`public interface UserService` in the same package as `UserServiceImpl`).
   - *Go Idiom*: Producers export concrete structs. Consumers declare tiny 1–2 method interfaces (e.g. `io.Reader`, `CustomerFinder`) specifying *only what the caller needs*. Duck typing satisfies them automatically.
3. **Errors as Values**:
   - *Java Equivalent*: Checked/Unchecked Exceptions (`try { ... } catch (IOException e)`).
   - *Go Idiom*: Multiple return values `(Result, error)` inspected explicitly with `if err != nil`. Errors wrapped with `fmt.Errorf("context: %w", err)` and inspected with `errors.Is` and `errors.As`.
4. **Channel Pipelines & Fan-Out / Fan-In**:
   - *Java Equivalent*: `ExecutorService`, `ForkJoinPool`, `CompletableFuture.allOf()`.
   - *Go Idiom*: Goroutines communicating over typed channels, synchronized using `sync.WaitGroup` or `golang.org/x/sync/errgroup`.
5. **Thread-Safe Lazy Initialization (`sync.Once`)**:
   - *Java Equivalent*: Double-checked locking `volatile instance` or `enum Singleton { INSTANCE; }`.
   - *Go Idiom*: `sync.Once.Do(func() { ... })` provides thread-safe, race-free single execution with minimal synchronization overhead.

---

### ⚠️ Top 10 Go Antipatterns for Java Developers

| # | Antipattern | Java Mental Model | What Breaks in Go | Idiomatic Go Fix |
|---|---|---|---|---|
| **1** | **The Nil Interface Trap** | Returning `null` exception reference means no exception. | An interface is `(Type, Value)`. Returning a typed `*MyError(nil)` produces an interface with non-nil Type, so `err != nil` evaluates to **`true`**! | Always return the untyped literal `nil` for no-error: `return nil`. |
| **2** | **Passing Mutex by Value** | Objects are references on the heap; passing passes the reference. | Structs are value types. Passing a struct with `sync.Mutex` creates a separate copy of the lock state, breaking mutual exclusion. | Always use pointer receivers (`*Counter`) and pass pointers. `go vet` catches this with `copylocks`. |
| **3** | **Goroutine Leaks** | Worker threads are bounded by `ExecutorService` or daemon pools. | Goroutines blocked on unbuffered channels with no receiver are never garbage collected, leaking stack memory and handles forever. | Use buffered channels (`make(chan T, 1)`) for single-result workers, and pass `context.Context` for cancellation. |
| **4** | **Interface Pollution** | Every service must have an `IService` and `ServiceImpl` for Spring DI. | Producer interfaces couple packages, bloat method contracts, and defeat Go's implicit structural duck typing. | Return concrete structs from providers. Let consumers define small (1–2 method) interfaces only where needed. |
| **5** | **Pointers to Reference Types** | Pass-by-value copies objects, so use pointers (`*[]T`, `*map[K]V`) for collections. | Slices, maps, and channels are already lightweight 24-byte headers or internal pointers. Pointer-to-map adds useless indirection and awkward syntax. | Pass slices, maps, and channels directly by value. |
| **6** | **Panic as Flow Control** | Exceptions are thrown for business validation (`throw new IllegalArgumentException()`). | `panic()` stops goroutine execution, unwinds the stack, and crashes the process if unhandled. It is not an exception mechanism. | Errors are values. Return `(T, error)` for business failures; reserve `panic()` strictly for unrecoverable bugs (programmer error). |
| **7** | **Slice Memory Retention** | `new ArrayList<>(list.subList(0, 2))` vs older Java 6 `substring` leak. | Sub-slicing `hugeArray[:4]` creates a slice header pointing to the original backing array, pinning megabytes/gigabytes in RAM. | Allocate a new small slice and copy only the needed elements: `small := make([]T, 4); copy(small, huge[:4])`. |
| **8** | **Slice Growth Without Preallocation** | `new ArrayList<>()` resizes automatically with negligible relative overhead in JVM. | Appending in a loop without capacity triggers repeated heap reallocations, memory copies, and GC pressure. | Preallocate capacity when known: `make([]T, 0, capacity)`. |
| **9** | **Unsynchronized Map Access** | `HashMap` isn't thread-safe, but concurrent access rarely crashes the JVM process. | The Go runtime detects concurrent map read/write and immediately terminates the process with a fatal, unrecoverable crash! | Synchronize with `sync.RWMutex`, `sync.Mutex`, or use `sync.Map`. |
| **10** | **Variable Shadowing of `err`** | Java compiler rejects duplicate local variable declarations in same scope. | The `:=` operator inside `if` or `for` declares a *new* local `err`, leaving the outer return `err` untouched and `nil`. | Use `=` assignment when setting existing outer variables, or structure error propagation explicitly. |

---

## 🧪 Running Tests

Go has a built-in testing framework.
- Run all tests: `go test -v ./...`
- Run a specific test: `go test -v ./pkg/generics -run TestMapValues`
- Run benchmarks: `go test -bench=. ./pkg/...`

---
*Happy Gophering! If you're stuck, look for the `// For a Java developer:` comments in the source code.*
