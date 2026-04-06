# AGENTS.md

## What this repo is
- `gogo` is a single-module Go 1.23 project (`go.mod`) that showcases language features, not a layered production service.
- The executable entrypoint is `cmd/gogo/main.go`; reusable demos live in sub-packages within `pkg/`.

## Big-picture architecture
- `cmd/gogo/main.go` is an ordered orchestrator that calls each demo in sequence, ending with a web server:
  `concurrency.RunConcurrencyDemo` → `interfaces.RunInterfacesDemo` → `generics.RunGenericsDemo` → `errors.RunErrorsDemo` → `collections.RunCollectionsDemo` → `basictypes.RunBasicTypesDemo` → `patterns.RunPatternsDemo` → `advanced.RunAdvancedDemo` → `iosystem.RunIOSystemDemo` → `commonlibs.RunCommonLibsDemo` → `web.StartWebServer("8080")`.
- Most demos are print-driven and side-effect-based (stdout), so behavior is observed by running the binary.

### File responsibilities
| Package/File | Purpose |
|---|---|
| `pkg/concurrency/` | Goroutines, channels, `sync.WaitGroup`, `context.WithTimeout`, Worker Pool, Context Value Propagation |
| `pkg/interfaces/` | `Shape` interface, `Rectangle`/`Circle`, polymorphism, type assertion, type switch |
| `pkg/generics/` | `List[T]`, `MapValues[T,R]` generic helpers |
| `pkg/errors/` | Custom errors (`MyCustomError`), `errors.As` |
| `pkg/collections/` | Generic data structures (`Set`, `Stack`, `Queue`, `OrderedMap`) and functional slice operations (`Filter`, `Reduce`, `FlatMap`, `GroupBy`, `Partition`, `Sorted`, `Distinct`, `Any`, `All`, `Zip`) – comparable to Java Collections + Stream API |
| `pkg/basictypes/` | Basic Go types manipulation: slices (lists), maps, strings, and JSON |
| `pkg/patterns/` | GoF design patterns: Singleton, Factory Method, Builder, Strategy, Observer, Decorator, Iterator, Adapter, Template Method, Command |
| `pkg/advanced/` | Advanced features: embedding (≈ inheritance), iota enums, bitmasks, closures/currying/memoization, defer/panic/recover, reflection, type constraints (`Number`), concurrent fan-out, struct tags/JSON, `embed` static assets. Also contains benchmarks, fuzzing, and Build System (Go vs Maven) comparison. |
| `pkg/iosystem/` | I/O and System-level programming: file manipulation, directory walking, env vars, CLI flags, exec commands, networking, signal handling |
| `pkg/commonlibs/` | Popular 3rd party libraries: `testify` (assertions), `zap` (logging), `uuid` (unique IDs), `gin` (web framework) |
| `pkg/web/` | HTTP handlers (`/hello`, `/json`, `/echo`) with logging middleware and JSON body parsing |

## Go Build System (for Java Developers)
- **Go tool**: A single executable CLI (`go`) handles build, test, package management, and more. It is a "batteries-included" tool, meaning it doesn't need plugins for core tasks (fmt, vet, test, build).
- **Go Modules (`go.mod`)**: Equivalent to Maven's `pom.xml` or Gradle's `build.gradle`. It manages dependencies and versioning in a simple text format.
- **`go build`**: Compiles the source code and dependencies into a single statically-linked binary (similar to a 'Fat JAR' but no JVM needed).
- **`go test`**: Runs built-in testing framework (similar to JUnit/Maven test). Includes support for benchmarks and fuzzing.
- **`go mod tidy`**: Cleans up and downloads missing dependencies (similar to Maven import/refresh).
- **`Makefile`**: Commonly used in Go projects to orchestrate CLI commands into named targets (similar to Maven lifecycle phases like `mvn package` or `mvn test`).
- **Deployment**: Java requires JRE/JDK on the target machine. Go produces a self-contained binary that can be copied and run without any pre-installed runtime.

## Critical workflows
- Run with Makefile (Recommended): `make run`
- Run full showcase: `go run ./cmd/gogo` (web server remains active until Enter is pressed in stdin).
- Run all tests: `make test` or `go test ./...`
- Build binary: `make build` (creates `gogo_binary`)
- Clean artifacts: `make clean`
- Run benchmarks: `go test -bench=. ./pkg/...`
- Run fuzz test: `go test -fuzz=FuzzReverse ./pkg/advanced`
- Validate web endpoints while app is running:
  - `curl "http://localhost:8080/hello?name=GoExpert"`
  - `curl http://localhost:8080/json`
  - `curl -X POST -H 'Content-Type: application/json' -d '{"msg":"hello"}' http://localhost:8080/echo`
- Focus a single test when iterating: `go test ./pkg/generics -run TestMapValues -v`.

## Project-specific coding patterns
- Keep feature demos discoverable by exposing top-level functions named `Run*Demo` in `pkg/*/`.
- Preserve zero-value usability where present (example: `Stack[T]{}`, `Queue[T]{}`, `List[T]{}` are valid before calling methods).
- Match existing testing style: direct `t.Errorf` checks, `reflect.DeepEqual` for slices, `math.Abs(…) > 1e-9` for float comparisons.
- Prefer standard library APIs already used in demos before introducing dependencies.
- GoF patterns use idiomatic Go: interfaces instead of abstract classes, embedding instead of inheritance, `sync.Once` for singletons, channels for iterators.

## Integration points and caveats
- No external services or third-party libs; integrations are Go stdlib (`net/http`, `encoding/json`, `context`, `sync`, `reflect`).
- `StartWebServer` uses global `http.HandleFunc`/default mux; registering the same routes twice in one process will fail.
- `GetAppConfig()` is a true process-wide singleton via `sync.Once` — calling it in tests shares state across the test binary.
- Server lifecycle is intentionally lightweight: launched with `go` routine, readiness is a fixed `100ms` sleep, and shutdown is process-driven.
- If adding new demos, wire them in `main.go` in sequence and keep output clearly sectioned like existing `--- Demo ---` / `--- Demo End ---` markers.
